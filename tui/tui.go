package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/awwwkshay/alpha-aide/agent/agent"
	"github.com/awwwkshay/alpha-aide/agent/config"
	agentmodels "github.com/awwwkshay/alpha-aide/agent/models"
	"github.com/awwwkshay/alpha-aide/agent/tools"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 22

// ── Slash commands ────────────────────────────────────────────────────────────

type SlashCmd struct {
	Name string
	Desc string
}

var slashCommands = []SlashCmd{
	{"/clear",    "Clear current session history"},
	{"/new",      "Open a new session"},
	{"/sessions", "Browse and switch sessions"},
	{"/models",   "Browse and switch models"},
	{"/help",     "Show keyboard shortcuts"},
	{"/exit",     "Quit alpha-aide"},
}

// ── Panel focus ───────────────────────────────────────────────────────────────

type panelFocus int

const (
	focusInput    panelFocus = iota
	focusSessions            // sidebar: sessions list
	focusModels              // sidebar: models list
)

// ── Tea messages ──────────────────────────────────────────────────────────────

type chunkMsg struct {
	chunk agent.StreamChunk
	ch    <-chan agent.StreamChunk
}

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	cfg      *config.Config
	allTools []agent.Tool

	width, height int
	focus         panelFocus

	// sessions
	sessions      []*agent.Session
	activeSession int
	sessionCursor int

	// models
	modelCursor int

	// input
	input textinput.Model

	// chat viewport
	viewport viewport.Model

	// slash autocomplete
	suggCursor int

	// streaming
	streaming bool
	cancelFn  context.CancelFunc

	status string
}

func initialModel(cfg *config.Config) (Model, error) {
	p, err := agent.NewProvider(cfg, cfg.Provider, cfg.Model)
	if err != nil {
		return Model{}, err
	}

	allTools := []agent.Tool{
		tools.ReadFileTool{},
		tools.WriteFileTool{},
		tools.EditFileTool{},
		tools.BashTool{},
	}

	a := agent.New(p, allTools)
	sess := agent.NewSession("main", cfg.Provider, cfg.Model, a)

	ti := textinput.New()
	ti.Placeholder = "Message or /command…"
	ti.Focus()
	ti.CharLimit = 4096

	vp := viewport.New(80, 20)
	vp.SetContent("")

	return Model{
		cfg:      cfg,
		allTools: allTools,
		sessions: []*agent.Session{sess},
		input:    ti,
		viewport: vp,
	}, nil
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeViewport()
		m.refreshViewport()
		return m, nil

	case chunkMsg:
		sess := m.sessions[m.activeSession]
		sess.HandleChunk(msg.chunk)
		m.refreshViewport()
		if msg.chunk.Type == "done" || msg.chunk.Type == "error" {
			m.streaming = false
			m.cancelFn = nil
			m.input.Focus()
			return m, nil
		}
		return m, nextChunk(msg.ch)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.focus == focusInput && !m.streaming {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleKey is the single key dispatch point.
func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	// ── Quit / cancel (always works) ─────────────────────────────────────
	if msg.Type == tea.KeyCtrlC {
		if m.streaming && m.cancelFn != nil {
			m.cancelFn()
			m.streaming = false
			m.cancelFn = nil
			m.input.Focus()
			m.focus = focusInput
			m.status = "cancelled"
			return m, nil
		}
		return m, tea.Quit
	}

	// ── Esc: always returns to input ──────────────────────────────────────
	if msg.Type == tea.KeyEsc {
		if strings.HasPrefix(m.input.Value(), "/") {
			m.input.SetValue("")
		}
		m.focus = focusInput
		m.input.Focus()
		m.status = ""
		m.resizeViewport()
		return m, nil
	}

	// ── Tab cycles focus: input → sessions → models → input ───────────────
	// Shift+Tab goes backwards.
	// (Tab does NOT trigger when suggestions are visible — handled below.)
	if msg.Type == tea.KeyTab || key == "shift+tab" {
		sugg := m.filteredSugg()
		if len(sugg) > 0 && msg.Type == tea.KeyTab {
			// Complete the highlighted suggestion instead of cycling focus
			m.input.SetValue(sugg[m.suggCursor].Name)
			m.input.CursorEnd()
			m.resizeViewport()
			return m, nil
		}
		if !m.streaming {
			backward := key == "shift+tab"
			m = m.cycleFocus(backward)
		}
		return m, nil
	}

	// ── Ctrl+N: new session ───────────────────────────────────────────────
	if key == "ctrl+n" {
		return m.newSession()
	}

	// ── Ctrl+W: close session ─────────────────────────────────────────────
	if key == "ctrl+w" {
		return m.deleteSession()
	}

	// ── Route to focused panel ────────────────────────────────────────────
	switch m.focus {
	case focusInput:
		return m.handleInputKey(msg)
	case focusSessions:
		return m.handleSessionsKey(msg)
	case focusModels:
		return m.handleModelsKey(msg)
	}
	return m, nil
}

// cycleFocus advances (or reverses) the panel focus.
func (m Model) cycleFocus(backward bool) Model {
	panels := []panelFocus{focusInput, focusSessions, focusModels}
	cur := 0
	for i, p := range panels {
		if p == m.focus {
			cur = i
			break
		}
	}
	if backward {
		cur = (cur - 1 + len(panels)) % len(panels)
	} else {
		cur = (cur + 1) % len(panels)
	}
	m.focus = panels[cur]
	if m.focus == focusInput {
		m.input.Focus()
		m.sessionCursor = m.activeSession
	} else {
		m.input.Blur()
		if m.focus == focusSessions {
			m.sessionCursor = m.activeSession
		}
	}
	return m
}

func (m Model) handleInputKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.streaming {
		// Allow viewport scroll while streaming
		switch msg.Type {
		case tea.KeyPgUp:
			m.viewport.HalfViewUp()
		case tea.KeyPgDown:
			m.viewport.HalfViewDown()
		}
		return m, nil
	}

	sugg := m.filteredSugg()

	switch msg.Type {

	case tea.KeyUp:
		if len(sugg) > 0 {
			if m.suggCursor > 0 {
				m.suggCursor--
			}
			return m, nil
		}
		// Scroll viewport when no suggestions
		m.viewport.HalfViewUp()
		return m, nil

	case tea.KeyDown:
		if len(sugg) > 0 {
			if m.suggCursor < len(sugg)-1 {
				m.suggCursor++
			}
			return m, nil
		}
		m.viewport.HalfViewDown()
		return m, nil

	case tea.KeyPgUp:
		m.viewport.HalfViewUp()
		return m, nil

	case tea.KeyPgDown:
		m.viewport.HalfViewDown()
		return m, nil

	case tea.KeyEnter:
		if len(sugg) > 0 {
			return m.executeSlash(sugg[m.suggCursor])
		}
		text := strings.TrimSpace(m.input.Value())
		if text == "" {
			return m, nil
		}
		if strings.HasPrefix(text, "/") {
			for _, c := range slashCommands {
				if c.Name == text {
					return m.executeSlash(c)
				}
			}
			m.status = "unknown command: " + text
			m.input.SetValue("")
			m.resizeViewport()
			return m, nil
		}
		m.input.SetValue("")
		m.suggCursor = 0
		m.resizeViewport()
		return m.submitMessage(text)
	}

	// All other keys — pass to textinput, then recalculate suggestions
	prevVal := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() != prevVal {
		m.suggCursor = 0
		m.resizeViewport()
	}
	return m, cmd
}

func (m Model) handleSessionsKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.sessionCursor > 0 {
			m.sessionCursor--
		}
	case tea.KeyDown:
		if m.sessionCursor < len(m.sessions)-1 {
			m.sessionCursor++
		}
	case tea.KeyEnter:
		m.activeSession = m.sessionCursor
		m.focus = focusInput
		m.input.Focus()
		m.refreshViewport()
	}
	return m, nil
}

func (m Model) handleModelsKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.modelCursor > 0 {
			m.modelCursor--
		}
	case tea.KeyDown:
		if m.modelCursor < len(KnownModels)-1 {
			m.modelCursor++
		}
	case tea.KeyEnter:
		sel := KnownModels[m.modelCursor]
		m.focus = focusInput
		m.input.Focus()
		return m.applyModelChange(sel)
	}
	return m, nil
}

// ── Slash command execution ───────────────────────────────────────────────────

func (m Model) executeSlash(cmd SlashCmd) (Model, tea.Cmd) {
	m.input.SetValue("")
	m.suggCursor = 0
	m.resizeViewport()

	switch cmd.Name {
	case "/clear":
		sess := m.sessions[m.activeSession]
		sess.Agent.ClearHistory()
		sess.Messages = nil
		sess.PendingText = ""
		sess.Turns = 0
		m.refreshViewport()
		m.status = "history cleared"

	case "/new":
		return m.newSession()

	case "/sessions":
		m.focus = focusSessions
		m.input.Blur()
		m.sessionCursor = m.activeSession

	case "/models":
		m.focus = focusModels
		m.input.Blur()
		// Sync cursor to current model
		sess := m.sessions[m.activeSession]
		for i, md := range KnownModels {
			if md.ID == sess.ModelID {
				m.modelCursor = i
				break
			}
		}

	case "/help":
		sess := m.sessions[m.activeSession]
		sess.Messages = append(sess.Messages, agent.DisplayMsg{
			Role: "info",
			Content: strings.TrimSpace(`Keyboard shortcuts:
  Tab / Shift+Tab   Cycle panel focus (input → sessions → models)
  ↑ / ↓             Navigate list or suggestions; scroll chat
  Enter             Select / send
  Tab               Complete slash command (when suggestion visible)
  Ctrl+N            New session
  Ctrl+W            Close session
  Esc               Return to input
  Ctrl+C            Cancel streaming / quit

Slash commands:
  /clear      Clear session history
  /new        New session
  /sessions   Browse sessions
  /models     Browse models
  /help       Show this help
  /exit       Quit`),
		})
		m.refreshViewport()

	case "/exit":
		return m, tea.Quit
	}

	return m, nil
}

// ── Actions ───────────────────────────────────────────────────────────────────

func (m Model) submitMessage(text string) (Model, tea.Cmd) {
	sess := m.sessions[m.activeSession]
	sess.Messages = append(sess.Messages, agent.DisplayMsg{Role: "user", Content: text})
	m.refreshViewport()

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFn = cancel
	m.streaming = true
	m.status = ""

	ch := make(chan agent.StreamChunk, 64)
	go func() {
		defer close(ch)
		err := sess.Agent.Run(ctx, text, func(chunk agent.StreamChunk) {
			select {
			case ch <- chunk:
			case <-ctx.Done():
			}
		})
		if err != nil && ctx.Err() == nil {
			ch <- agent.StreamChunk{Type: "error", Err: err}
		}
	}()

	return m, nextChunk(ch)
}

func (m Model) newSession() (Model, tea.Cmd) {
	n := len(m.sessions) + 1
	sess := m.sessions[m.activeSession]
	p, err := agent.NewProvider(m.cfg, sess.ProviderName, sess.ModelID)
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	a := agent.New(p, m.allTools)
	newSess := agent.NewSession(fmt.Sprintf("session-%d", n), sess.ProviderName, sess.ModelID, a)
	m.sessions = append(m.sessions, newSess)
	m.activeSession = len(m.sessions) - 1
	m.sessionCursor = m.activeSession
	m.focus = focusInput
	m.input.Focus()
	m.input.SetValue("")
	m.refreshViewport()
	return m, nil
}

func (m Model) deleteSession() (Model, tea.Cmd) {
	if len(m.sessions) <= 1 {
		m.status = "cannot close last session"
		return m, nil
	}
	m.sessions = append(m.sessions[:m.activeSession], m.sessions[m.activeSession+1:]...)
	if m.activeSession >= len(m.sessions) {
		m.activeSession = len(m.sessions) - 1
	}
	m.sessionCursor = m.activeSession
	m.focus = focusInput
	m.input.Focus()
	m.refreshViewport()
	return m, nil
}

func (m Model) applyModelChange(sel ModelDef) (Model, tea.Cmd) {
	sess := m.sessions[m.activeSession]
	p, err := agent.NewProvider(m.cfg, sel.Provider, sel.ID)
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	newAgent := agent.NewWithHistory(p, m.allTools, sess.Agent.History())
	sess.Agent = newAgent
	sess.ProviderName = sel.Provider
	sess.ModelID = sel.ID
	m.status = fmt.Sprintf("switched to %s", sel.Display)
	return m, nil
}

func nextChunk(ch <-chan agent.StreamChunk) tea.Cmd {
	return func() tea.Msg {
		chunk, ok := <-ch
		if !ok {
			return chunkMsg{chunk: agent.StreamChunk{Type: "done"}, ch: nil}
		}
		return chunkMsg{chunk: chunk, ch: ch}
	}
}

// ── Suggestion helpers ────────────────────────────────────────────────────────

func (m Model) filteredSugg() []SlashCmd {
	val := m.input.Value()
	if !strings.HasPrefix(val, "/") {
		return nil
	}
	var out []SlashCmd
	for _, c := range slashCommands {
		if strings.HasPrefix(c.Name, val) {
			out = append(out, c)
		}
	}
	return out
}

func (m Model) suggBoxHeight() int {
	n := len(m.filteredSugg())
	if n == 0 {
		return 0
	}
	if n > 6 {
		n = 6
	}
	return n + 2 // top + bottom border
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.width == 0 {
		return "Starting…"
	}
	parts := []string{
		m.renderHeader(),
		lipgloss.JoinHorizontal(lipgloss.Top, m.renderSidebar(), m.renderChat()),
	}
	if sugg := m.renderSuggestions(); sugg != "" {
		parts = append(parts, sugg)
	}
	parts = append(parts, m.renderInput(), m.renderHelp())
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// ── Render helpers ────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	sess := m.sessions[m.activeSession]
	logo := logoStyle.Render("alpha-aide")
	modelLabel := headerStyle.Render(agentmodels.ShortName(sess.ModelID))
	stats := dimStyle.Render(sess.StatsLine())

	var tabs []string
	for i, s := range m.sessions {
		if i == m.activeSession {
			tabs = append(tabs, activeItemStyle.Render("[ "+s.Name+" ]"))
		} else {
			tabs = append(tabs, dimStyle.Render("  "+s.Name+"  "))
		}
	}

	left := logo + "  " + strings.Join(tabs, "")
	right := modelLabel + "  " + stats
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#111111")).
		Width(m.width).
		Render(left + strings.Repeat(" ", gap) + right)
}

func (m Model) renderSidebar() string {
	h := m.height - 3 - m.suggBoxHeight()

	var sb strings.Builder

	// ── Sessions ──
	if m.focus == focusSessions {
		sb.WriteString(activeItemStyle.Render("▶ Sessions") + "\n")
	} else {
		sb.WriteString(sectionStyle.Render("  Sessions") + "\n")
	}
	for i, s := range m.sessions {
		label := truncatePad(s.Name, sidebarWidth-4)
		prefix := "  "
		if i == m.activeSession {
			prefix = cursorStyle.Render("▸ ")
		}
		row := prefix + label
		highlighted := m.focus == focusSessions && i == m.sessionCursor
		if highlighted || i == m.activeSession {
			sb.WriteString(activeItemStyle.Render(row) + "\n")
		} else {
			sb.WriteString(inactiveItemStyle.Render(row) + "\n")
		}
	}
	sb.WriteString("\n")

	// ── Models ──
	if m.focus == focusModels {
		sb.WriteString(activeItemStyle.Render("▶ Models") + "\n")
	} else {
		sb.WriteString(sectionStyle.Render("  Models") + "\n")
	}
	sess := m.sessions[m.activeSession]
	for i, md := range KnownModels {
		label := truncatePad(md.Display, sidebarWidth-4)
		prefix := "  "
		if md.ID == sess.ModelID {
			prefix = cursorStyle.Render("▸ ")
		}
		row := prefix + label
		highlighted := m.focus == focusModels && i == m.modelCursor
		if highlighted || md.ID == sess.ModelID {
			sb.WriteString(activeItemStyle.Render(row) + "\n")
		} else {
			sb.WriteString(inactiveItemStyle.Render(row) + "\n")
		}
	}

	style := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(h).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(dim).
		PaddingLeft(1)
	return style.Render(sb.String())
}

func (m Model) renderChat() string {
	return m.viewport.View()
}

func (m Model) renderSuggestions() string {
	sugg := m.filteredSugg()
	if len(sugg) == 0 {
		return ""
	}

	maxVisible := 6
	start := 0
	if m.suggCursor >= maxVisible {
		start = m.suggCursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(sugg) {
		end = len(sugg)
	}
	visible := sugg[start:end]

	nameColW := 11
	innerW := m.width - 4

	var rows []string
	for i, c := range visible {
		actualIdx := start + i
		name := lipgloss.NewStyle().Width(nameColW).Render(c.Name)
		desc := c.Desc
		maxDesc := innerW - nameColW - 2
		if maxDesc > 0 && len(desc) > maxDesc {
			desc = desc[:maxDesc-1] + "…"
		}
		row := name + "  " + desc
		if actualIdx == m.suggCursor {
			rows = append(rows, activeItemStyle.Render(" "+row))
		} else {
			rows = append(rows, inactiveItemStyle.Render(" "+row))
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cyan).
		Width(m.width - 2).
		Render(strings.Join(rows, "\n"))
}

func (m Model) renderInput() string {
	prefix := inputPrefixStyle.Render("> ")
	suffix := ""
	if m.streaming {
		suffix = dimStyle.Render("  [streaming…]")
	}
	if m.status != "" {
		suffix += "  " + dimStyle.Render(m.status)
	}
	return lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(dim).
		Render(prefix + m.input.View() + suffix)
}

func (m Model) renderHelp() string {
	var hint string
	switch m.focus {
	case focusSessions:
		hint = "  ↑↓ navigate  ·  Enter select  ·  Tab next panel  ·  Esc back"
	case focusModels:
		hint = "  ↑↓ navigate  ·  Enter switch model  ·  Tab next panel  ·  Esc back"
	default:
		hint = "  /command  ·  Tab sessions  ·  Ctrl+N new  ·  Ctrl+W close  ·  ↑↓ scroll  ·  Ctrl+C quit"
	}
	return helpStyle.Render(hint)
}

// ── Viewport helpers ──────────────────────────────────────────────────────────

func (m *Model) refreshViewport() {
	sess := m.sessions[m.activeSession]
	m.viewport.SetContent(renderMessages(sess))
	m.viewport.GotoBottom()
}

func (m *Model) resizeViewport() {
	chatW := m.width - sidebarWidth - 2
	chatH := m.height - 3 - m.suggBoxHeight()
	if chatW < 10 {
		chatW = 10
	}
	if chatH < 3 {
		chatH = 3
	}
	m.viewport.Width = chatW
	m.viewport.Height = chatH
}

func renderMessages(sess *agent.Session) string {
	var sb strings.Builder
	for _, msg := range sess.Messages {
		switch msg.Role {
		case "user":
			sb.WriteString(userLabelStyle.Render("You") + "\n")
			sb.WriteString(msg.Content + "\n\n")
		case "assistant":
			sb.WriteString(agentLabelStyle.Render("Agent") + "\n")
			sb.WriteString(msg.Content + "\n\n")
		case "tool_call":
			sb.WriteString(toolStyle.Render(fmt.Sprintf("[tool: %s]", msg.Tool)) + "\n\n")
		case "tool_result":
			sb.WriteString(toolStyle.Render(fmt.Sprintf("[result: %s]", msg.Tool)) + "\n")
			sb.WriteString(toolResultStyle.Render(msg.Content) + "\n\n")
		case "info":
			sb.WriteString(dimStyle.Render(msg.Content) + "\n\n")
		case "system":
			sb.WriteString(errorStyle.Render(msg.Content) + "\n\n")
		}
	}
	if sess.PendingText != "" {
		sb.WriteString(agentLabelStyle.Render("Agent") + "\n")
		sb.WriteString(sess.PendingText)
	}
	return sb.String()
}

// ── Utilities ─────────────────────────────────────────────────────────────────

func truncatePad(s string, n int) string {
	if len(s) > n {
		return s[:n-1] + "…"
	}
	return s
}

