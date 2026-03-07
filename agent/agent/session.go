package agent

import (
	"fmt"

	"github.com/awwwkshay/alpha-aide/llm"
)

// DisplayMsg is a rendered message for display.
type DisplayMsg struct {
	Role    string // "user" | "assistant" | "tool_call" | "tool_result" | "system" | "info"
	Content string
	Tool    string // tool name for tool_call / tool_result
}

// Session holds one conversation thread.
type Session struct {
	Name         string
	Agent        *Agent
	ProviderName string
	ModelID      string
	Messages     []DisplayMsg
	PendingText  string // accumulates streaming assistant text
	Turns        int    // completed turns
}

func NewSession(name, providerName, modelID string, a *Agent) *Session {
	return &Session{
		Name:         name,
		Agent:        a,
		ProviderName: providerName,
		ModelID:      modelID,
	}
}

func (s *Session) StatsLine() string {
	if s.Turns == 0 {
		return ""
	}
	return fmt.Sprintf("%d turn(s)", s.Turns)
}

// HandleChunk processes one streaming chunk into the session's message list.
func (s *Session) HandleChunk(chunk llm.StreamChunk) {
	switch chunk.Type {
	case "text_delta":
		s.PendingText += chunk.Text

	case "tool_use_start":
		if s.PendingText != "" {
			s.Messages = append(s.Messages, DisplayMsg{Role: "assistant", Content: s.PendingText})
			s.PendingText = ""
		}
		s.Messages = append(s.Messages, DisplayMsg{
			Role: "tool_call",
			Tool: chunk.ToolCall.Name,
		})

	case "tool_result":
		content := chunk.Text
		if len(content) > 300 {
			content = content[:300] + "…"
		}
		s.Messages = append(s.Messages, DisplayMsg{
			Role:    "tool_result",
			Tool:    chunk.ToolCall.Name,
			Content: content,
		})

	case "done":
		if s.PendingText != "" {
			s.Messages = append(s.Messages, DisplayMsg{Role: "assistant", Content: s.PendingText})
			s.PendingText = ""
		}
		s.Turns++

	case "error":
		if s.PendingText != "" {
			s.Messages = append(s.Messages, DisplayMsg{Role: "assistant", Content: s.PendingText})
			s.PendingText = ""
		}
		s.Messages = append(s.Messages, DisplayMsg{
			Role:    "system",
			Content: "Error: " + chunk.Err.Error(),
		})
	}
}
