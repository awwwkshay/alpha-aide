# alpha-aide-tui

A terminal UI for the alpha-aide coding agent. Manage multiple sessions, switch models, and stream responses — all from the terminal.

## Layout

```
┌─ alpha-aide  [ main ]  [ session-2 ] ────────────── haiku-4.5 ─┐
├─ Sessions ──────┬──────────────────────────────────────────────┤
│   Sessions      │  You                                         │
│ ▸ main          │  write a hello world in go                   │
│   session-2     │                                              │
│                 │  Agent                                        │
│   Models        │  Here's a simple hello world program:        │
│ ▸ Claude Haiku  │  ...                                         │
│   Claude Sonnet │                                              │
│   GPT-4o        │  [tool: bash]                                │
│   Gemini Flash  │                                              │
│   …             │  [result: bash]                              │
│                 │  Hello, World!                               │
├─────────────────┴──────────────────────────────────────────────┤
│ > _                                                            │
│  ^N new  ·  ^W close  ·  ^S sessions  ·  ^M models  ·  ^C quit│
└────────────────────────────────────────────────────────────────┘
```

## Running

```bash
# With Makefile (recommended)
make run-tui                         # uses env vars
make run-tui-anthropic               # Claude Haiku 4.5
make run-tui-gemini                  # Gemini 2.0 Flash
make run-tui-openai                  # GPT-4o Mini

# Direct binary (after make build-tui)
./dist/alpha-aide-tui
./dist/alpha-aide-tui -provider anthropic -model claude-sonnet-4-6

# Without building
ANTHROPIC_API_KEY=sk-ant-... go run ./tui
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Ctrl+N` | New session (inherits current model) |
| `Ctrl+W` | Close current session |
| `Ctrl+S` | Focus sessions panel — navigate with `↑↓`, confirm with `Enter` |
| `Ctrl+M` | Focus models panel — navigate with `↑↓`, confirm with `Enter` |
| `Esc` | Return focus to input |
| `Ctrl+C` | Cancel streaming request (if active), or quit |

## Sessions

Each session is an independent conversation with its own history. New sessions inherit the provider and model of the current session. You can have unlimited sessions open at once and switch between them instantly.

## Model Switching

Press `Ctrl+M` to open the models panel. Navigate with arrow keys and press `Enter` to switch. The switch applies to the current session and preserves its conversation history — the next message simply uses the new model.

Available models (from `agent/models`):

| Provider | Models |
|----------|--------|
| Anthropic | Claude Opus 4.6, Claude Sonnet 4.6, Claude Haiku 4.5 |
| OpenAI | GPT-4o, GPT-4o Mini, o3 Mini |
| Gemini | Gemini 2.0 Flash, Gemini 1.5 Pro, Gemini 1.5 Flash |

## Environment Variables

Same as the agent CLI — the TUI reads from `.env` or the environment:

| Variable | Description |
|----------|-------------|
| `AGENT_PROVIDER` | `anthropic` (default), `openai`, `gemini` |
| `AGENT_MODEL` | Model ID (e.g. `claude-haiku-4-5-20251001`) |
| `ANTHROPIC_API_KEY` | Required for Anthropic |
| `OPENAI_API_KEY` | Required for OpenAI |
| `OPENAI_BASE_URL` | Override for Groq, Ollama, etc. |
| `GEMINI_API_KEY` | Required for Gemini |
| `AGENT_MAX_TOKENS` | Max response tokens (default 8096) |

## Architecture

The TUI is a separate Go module (`github.com/awwwkshay/alpha-aide/tui`) that imports the agent SDK. It adds no dependencies to the agent or llm modules.

```
tui/
├── main.go      Entry point — loads config, starts Bubble Tea program
├── tui.go       Bubble Tea model: Init / Update / View
├── session.go   Session type: history, streaming chunk handler
├── models.go    Re-exports model catalog from agent/models
└── styles.go    Lipgloss color theme
```

**Key libraries:**
- [`charmbracelet/bubbletea`](https://github.com/charmbracelet/bubbletea) — Elm-architecture TUI framework
- [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) — Layout and styling
- [`charmbracelet/bubbles`](https://github.com/charmbracelet/bubbles) — Text input and viewport components
