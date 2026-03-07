package llm

import "context"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type ContentBlock struct {
	Type       string
	Text       string
	ToolUseID  string
	ToolName   string
	Content    string
	ToolInput  map[string]any
	IsError    bool
}

type Message struct {
	Role    Role
	Content []ContentBlock
}

type ToolCall struct {
	ID    string
	Name  string
	Input map[string]any
}

type ToolResult struct {
	ToolUseID string
	Content   string
	IsError   bool
}

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]any
	Execute(ctx context.Context, input map[string]any) ToolResult
}

type StreamChunk struct {
	Type     string    // "text_delta" | "tool_use_start" | "tool_use_end" | "error" | "done"
	Text     string
	ToolCall *ToolCall // non-nil on "tool_use_start"
	Err      error
}

type Provider interface {
	Complete(ctx context.Context, systemPrompt string, messages []Message, tools []Tool) (<-chan StreamChunk, error)
}
