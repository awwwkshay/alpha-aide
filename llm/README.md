# llm

A minimal, multi-provider streaming LLM library for Go.

**Module:** `github.com/awwwkshay/alpha-aide/llm`

## Providers

| Provider          | Package          | Notes                                              |
| ----------------- | ---------------- | -------------------------------------------------- |
| Anthropic         | `llm/anthropic`  | Native SDK, streaming                              |
| OpenAI-compatible | `llm/openai`     | Works with Groq, Ollama, any OpenAI-compatible API |
| Gemini            | `llm/gemini`     | Google Gemini via `google.golang.org/genai`        |

## Core Types

Defined in `llm.go`:

```go
// Implement this to add a new provider
type Provider interface {
    Complete(ctx context.Context, systemPrompt string, messages []Message, tools []Tool) (<-chan StreamChunk, error)
}

// Implement this to add a new tool
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]any
    Execute(ctx context.Context, input map[string]any) ToolResult
}

// Stream chunks from Complete()
type StreamChunk struct {
    Type     string    // "text_delta" | "tool_use_start" | "tool_use_end" | "error" | "done"
    Text     string
    ToolCall *ToolCall // non-nil on "tool_use_start"
    Err      error
}
```

## Usage

### Anthropic

```go
import "github.com/awwwkshay/alpha-aide/llm/anthropic"

provider := anthropic.New(os.Getenv("ANTHROPIC_API_KEY"), "claude-haiku-4-5-20251001", 8096)
```

### OpenAI-compatible

```go
import "github.com/awwwkshay/alpha-aide/llm/openai"

// OpenAI
provider := openai.New(os.Getenv("OPENAI_API_KEY"), "https://api.openai.com/v1", "gpt-4o", 8096)

// Groq
provider := openai.New(os.Getenv("OPENAI_API_KEY"), "https://api.groq.com/openai/v1", "llama-3.1-70b-versatile", 8096)

// Ollama (local)
provider := openai.New("", "http://localhost:11434/v1", "qwen2.5-coder:32b", 8096)
```

### Gemini

```go
import "github.com/awwwkshay/alpha-aide/llm/gemini"

provider := gemini.New(os.Getenv("GEMINI_API_KEY"), "gemini-2.0-flash", 8096)
```

### Streaming

```go
ch, err := provider.Complete(ctx, systemPrompt, messages, tools)
if err != nil {
    return err
}
for chunk := range ch {
    switch chunk.Type {
    case "text_delta":
        fmt.Print(chunk.Text)
    case "tool_use_start":
        fmt.Printf("\n[calling %s]\n", chunk.ToolCall.Name)
    case "error":
        return chunk.Err
    case "done":
        // stream finished
    }
}
```

## Dependencies

- `github.com/anthropics/anthropic-sdk-go v0.2.0-alpha.13`
- `github.com/openai/openai-go v0.1.0-alpha.62`
- `google.golang.org/genai v1.49.0`
