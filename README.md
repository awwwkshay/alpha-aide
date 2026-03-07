# alpha-aide

A minimal coding agent in Go. Inspired by the philosophy of radical minimalism — 4 tools, no permission popups, shortest viable system prompt.

Two interfaces: a plain REPL (`agent`) and a full terminal UI (`tui`).

## Structure

```text
alpha-aide/
├── go.work          # Go workspace linking all modules
├── llm/             # Reusable multi-provider LLM library
├── agent/           # Coding agent SDK + CLI
│   └── models/      # Known model catalog
└── tui/             # Terminal UI (sessions, model switcher)
```

## Quick Start

### REPL (plain)

```bash
ANTHROPIC_API_KEY=your-key go run ./agent
```

### TUI

```bash
ANTHROPIC_API_KEY=your-key go run ./tui
```

### With Make

```bash
make run            # agent REPL (Anthropic default)
make run-tui        # terminal UI (Anthropic default)
make run-tui-gemini # TUI with Gemini 2.0 Flash
make run-tui-openai # TUI with GPT-4o Mini
```

See [`docs/getting-started.md`](docs/getting-started.md) for full setup instructions.

## Tools

| Tool | Description |
| --- | --- |
| `read_file` | Read file contents |
| `write_file` | Write file, creates parent dirs |
| `edit_file` | Replace unique string in a file |
| `bash` | Run shell commands |

## REPL Commands

| Command | Action |
| --- | --- |
| `/help` | Show commands |
| `/tools` | List tools and descriptions |
| `/clear` | Clear conversation history |
| `/history` | Show conversation history |
| `/exit` | Exit |

## TUI Shortcuts

| Key | Action |
| --- | --- |
| `Ctrl+N` | New session |
| `Ctrl+W` | Close session |
| `Ctrl+S` | Switch session |
| `Ctrl+M` | Switch model |
| `Ctrl+C` | Cancel / quit |

## Modules

- [`llm/`](llm/README.md) — Multi-provider streaming LLM library
- [`agent/`](agent/README.md) — Coding agent SDK and CLI
- [`tui/`](tui/README.md) — Terminal UI for sessions and model management

## Docs

- [`docs/getting-started.md`](docs/getting-started.md) — Installation and first run
- [`docs/development.md`](docs/development.md) — Architecture and contribution guide
