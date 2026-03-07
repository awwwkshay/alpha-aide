# Development Guide

## Repo Layout

```text
alpha-aide/
├── go.work              Go workspace (links ./llm, ./agent, ./tui)
├── go.work.sum
├── Makefile             Build + run targets; outputs to ./dist/
├── docs/
│   ├── getting-started.md
│   └── development.md   (this file)
├── llm/                 Standalone LLM library module
│   ├── go.mod           module github.com/awwwkshay/alpha-aide/llm
│   ├── llm.go           Core interfaces: Provider, Tool, Message, StreamChunk
│   ├── anthropic/       Anthropic native provider
│   ├── openai/          OpenAI-compatible provider
│   └── gemini/          Google Gemini provider
├── agent/               Agent SDK + CLI module
│   ├── go.mod           module github.com/awwwkshay/alpha-aide/agent
│   ├── main.go          Entry point (plain REPL)
│   ├── agent/           ReAct loop
│   ├── models/          Known model catalog (shared with TUI)
│   ├── tools/           4 tool implementations
│   ├── repl/            REPL + colors
│   └── config/          Env + flag config
└── tui/                 Terminal UI module
    ├── go.mod           module github.com/awwwkshay/alpha-aide/tui
    ├── main.go          Entry point
    ├── tui.go           Bubble Tea model (Init/Update/View)
    ├── session.go       Session type + streaming chunk handler
    ├── models.go        Re-exports agent/models catalog
    ├── styles.go        Lipgloss color theme
    └── README.md        TUI keyboard reference
```

## Go Workspace Setup

This repo uses [Go workspaces](https://go.dev/ref/mod#workspaces) so all modules (`llm`, `agent`, `tui`) can reference each other locally without publishing.

**Important:** `go.work` alone is not sufficient. Each `go.mod` also has `replace` directives for its local dependencies:

```text
require github.com/awwwkshay/alpha-aide/llm v0.0.0-00010101000000-000000000000
replace github.com/awwwkshay/alpha-aide/llm => ../llm
```

The sentinel version `v0.0.0-00010101000000-000000000000` tells `go mod tidy` this module is never fetched from a proxy.

### After adding a new dependency

```bash
cd <module> && go mod tidy
```

### After changing `llm` interfaces

```bash
# Workspace handles local resolution — no version bumping needed
make build
```

### Adding a new module to the workspace

1. Create the module directory with `go mod init`
2. Add `replace` directives for any local deps in the new `go.mod`
3. Add `./newmodule` to `go.work` under `use`

## Adding a New Tool

1. Create `agent/tools/mytool.go`:

```go
package tools

import (
    "context"
    "github.com/awwwkshay/alpha-aide/llm"
)

type MyTool struct{}

func (MyTool) Name() string        { return "my_tool" }
func (MyTool) Description() string { return "What this tool does." }
func (MyTool) InputSchema() map[string]any {
    return map[string]any{
        "type": "object",
        "properties": map[string]any{
            "param": map[string]any{"type": "string", "description": "..."},
        },
        "required": []string{"param"},
    }
}

func (MyTool) Execute(_ context.Context, input map[string]any) llm.ToolResult {
    param, _ := input["param"].(string)
    // ... do work ...
    return llm.ToolResult{Content: "result"}
}
```

1. Register it in `agent/main.go`:

```go
allTools := []llm.Tool{
    tools.ReadFileTool{},
    tools.WriteFileTool{},
    tools.EditFileTool{},
    tools.BashTool{},
    tools.MyTool{},  // add here
}
```

1. Mention it in the system prompt if its usage needs guidance.

## Adding a New Provider

Implement `llm.Provider` in a new package under `llm/`:

```go
package myprovider

import (
    "context"
    "github.com/awwwkshay/alpha-aide/llm"
)

type Provider struct{ /* ... */ }

func New( /* config */ ) *Provider { return &Provider{} }

func (p *Provider) Complete(ctx context.Context, systemPrompt string, messages []llm.Message, tools []llm.Tool) (<-chan llm.StreamChunk, error) {
    ch := make(chan llm.StreamChunk, 64)
    go func() {
        defer close(ch)
        // stream chunks: "text_delta", "tool_use_start", "tool_use_end", "done", "error"
    }()
    return ch, nil
}
```

Then wire it in `agent/main.go` under a new provider name and add its API key to `config.go`.

## StreamChunk Protocol

The `Complete()` channel must emit chunks in this order:

```text
text_delta*           (zero or more text chunks)
(tool_use_start       (one per tool call)
 tool_use_end)*
done                  (always last, or)
error                 (on failure)
```

The agent accumulates all tool calls before executing any. Tool results are appended as a single `user` message with multiple `tool_result` blocks.

## Key Design Decisions

**Uniqueness invariant in `edit_file`:** `old_str` must appear exactly once. 0 occurrences → error. >1 occurrences → error with count. This prevents silent multi-replace bugs and forces the LLM to provide enough context.

**Non-zero bash exits are not errors:** The `bash` tool never sets `IsError: true` for non-zero exit codes. The LLM sees the full output (stdout + stderr + exit code) and decides what to do next. Only actual execution failures (timeout, binary not found) skip the output.

**SIGINT cancels, not exits:** Ctrl+C cancels the current `context.Context`, which stops the in-flight request. The process keeps running and returns to the `>` prompt.

**No session persistence:** History is in-memory only. Use `/clear` to reset, or rely on file-based statefulness (`TODO.md`, `PLAN.md`, etc.) for long tasks.

## Running Tests

```bash
# No unit tests yet — the codebase is verified via manual smoke tests
go build ./...   # from repo root via workspace

# Smoke test the REPL
printf '/help\n/tools\n/exit\n' | ANTHROPIC_API_KEY=test go run ./agent
```

## Building for Distribution

```bash
# Both binaries into dist/
make build

# Individual
make build-agent   # → dist/alpha-aide
make build-tui     # → dist/alpha-aide-tui

# Cross-compile (agent example)
cd agent
GOOS=linux  GOARCH=amd64 go build -o ../dist/alpha-aide-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o ../dist/alpha-aide-darwin-arm64 .
```

## Publishing `llm` as a Standalone Module

When ready to decouple `llm` from this monorepo:

1. Tag the release: `git tag llm/v0.1.0 && git push origin llm/v0.1.0`
2. Remove the `replace` directive from `agent/go.mod`
3. Update the `require` version: `cd agent && go get github.com/awwwkshay/alpha-aide/llm@v0.1.0`
4. The `go.work` `use ./llm` can remain for local development or be removed

No restructuring needed — the module boundary is already clean.
