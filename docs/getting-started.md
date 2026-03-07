# Getting Started

## Prerequisites

- Go 1.26 or later
- An API key for your chosen provider

## Installation

```bash
git clone https://github.com/awwwkshay/alpha-aide.git
cd alpha-aide
```

No build step required — run directly with `go run`.

## Running

### With Anthropic (default)

```bash
ANTHROPIC_API_KEY=sk-ant-... go run ./agent
```

### With OpenAI

```bash
AGENT_PROVIDER=openai \
OPENAI_API_KEY=sk-... \
AGENT_MODEL=gpt-4o \
go run ./agent
```

### With Groq (free tier available)

```bash
AGENT_PROVIDER=openai \
OPENAI_API_KEY=gsk_... \
AGENT_MODEL=llama-3.1-70b-versatile \
OPENAI_BASE_URL=https://api.groq.com/openai/v1 \
go run ./agent
```

### With Gemini

```bash
AGENT_PROVIDER=gemini \
GEMINI_API_KEY=your-key \
AGENT_MODEL=gemini-2.0-flash \
go run ./agent
```

### With Ollama (local, no API key)

```bash
# Start Ollama first: ollama serve
AGENT_PROVIDER=openai \
OPENAI_API_KEY=ollama \
AGENT_MODEL=qwen2.5-coder:32b \
OPENAI_BASE_URL=http://localhost:11434/v1 \
go run ./agent
```

## First Session

Once the REPL starts you'll see:

```text
Alpha Aide Agent — type /help for commands

>
```

Try these to verify everything works:

```text
> list files in the current directory
> read the go.work file
> create a file called /tmp/hello.txt with "hello world" in it
> run go build ./... and show me the output
```

## REPL Commands

| Command | Action |
| --- | --- |
| `/help` | Show all commands |
| `/tools` | List tools and descriptions |
| `/clear` | Clear conversation history |
| `/history` | Print conversation history |
| `/exit` or `/quit` | Exit |

**Ctrl+C** cancels the current request and returns to the `>` prompt without exiting.

## Tips

**Track complex tasks in a file.** For multi-step work, tell the agent:

```text
> keep a TODO.md to track your progress on this task
```

**Use bash to verify changes.** The agent will automatically run tests or build commands when you ask it to verify its work.

**Be specific about scope.** The agent defaults to targeted edits. If you want a full rewrite, say so explicitly.

## Using the TUI

The TUI provides a full-screen interface with multiple sessions and live model switching.

```bash
# Run directly
ANTHROPIC_API_KEY=sk-ant-... go run ./tui

# With a different provider
AGENT_PROVIDER=gemini GEMINI_API_KEY=your-key go run ./tui
```

See [`tui/README.md`](../tui/README.md) for the full keyboard reference.

## Building Binaries

```bash
make build           # builds both dist/alpha-aide and dist/alpha-aide-tui
make build-agent     # agent REPL only
make build-tui       # TUI only
```

Or build manually:

```bash
go build -C agent -o ../dist/alpha-aide .
go build -C tui   -o ../dist/alpha-aide-tui .
```
