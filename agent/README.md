# agent

The alpha-aide CLI. Wraps the `llm` library with 4 tools, a ReAct loop, and a simple REPL.

**Module:** `github.com/awwwkshay/alpha-aide/agent`

## Usage

```bash
# From the repo root
ANTHROPIC_API_KEY=your-key go run ./agent

# With flags
go run ./agent -provider anthropic -model claude-haiku-4-5-20251001 -max-tokens 8096
```

## Configuration

Priority: CLI flags > environment variables > defaults.

| Env Var | Flag | Default | Description |
| --------- | --------- | --------- | ------------- |
| `ANTHROPIC_API_KEY` | — | — | Anthropic API key |
| `OPENAI_API_KEY` | — | — | OpenAI-compatible API key |
| `GEMINI_API_KEY` | — | — | Google Gemini API key |
| `OPENAI_BASE_URL` | `-baseurl` | _(empty)_ | Base URL for OpenAI-compatible provider (defaults to `api.openai.com`) |
| `AGENT_PROVIDER` | `-provider` | `anthropic` | `anthropic`, `openai`, or `gemini` |
| `AGENT_MODEL` | `-model` | `claude-haiku-4-5-20251001` | Model name |
| `NO_COLOR` | `-no-color` | false | Disable ANSI colors |
| — | `-max-tokens` | `8096` | Max response tokens |

A `.env` file in the working directory (or any parent directory) is automatically loaded. Variables already set in the environment take precedence.

## Packages

```text
agent/
├── main.go            Entry point — parse flags, wire provider + tools, start REPL
├── agent/
│   ├── agent.go       ReAct loop: streams → accumulates → executes tools → repeats
│   └── system_prompt.go  ~110-token system prompt constant
├── tools/
│   ├── tool.go        Re-exports llm.Tool and llm.ToolResult
│   ├── read.go        read_file
│   ├── write.go       write_file
│   ├── edit.go        edit_file (uniqueness-enforced str_replace)
│   └── bash.go        bash (stdout/stderr, exit code, 100KB cap, UTF-8 sanitize)
├── repl/
│   ├── repl.go        REPL loop + slash command dispatch + SIGINT handling
│   └── colors.go      ANSI colors, auto-disabled when not a TTY
└── config/
    └── config.go      Config struct loaded from env + flags
```

## Tools

### `edit_file` — uniqueness invariant

`old_str` must appear **exactly once** in the file. If it appears 0 times you get an error; if it appears more than once you get a count with a prompt to add more context. This prevents accidental multi-replace bugs.

### `bash`

- Timeout: default 30s, max 300s (set via `timeout_seconds` param)
- Output: stdout and stderr each capped at 100 KB
- Non-zero exit codes are **not** treated as errors — the LLM sees the output and decides what to do next

## Agent Loop

```text
user input
    ↓
append to history
    ↓
loop (max 20 iterations):
    call provider.Complete() → stream chunks
    accumulate text + tool calls
    append assistant message to history
    if no tool calls → done
    execute each tool → collect results
    append tool_result message to history
```
