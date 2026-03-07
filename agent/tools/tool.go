package tools

import (
	"context"

	"github.com/awwwkshay/alpha-aide/llm"
)

// Re-export llm types for convenience within the agent package.
type Tool = llm.Tool
type ToolResult = llm.ToolResult

// Registry holds all registered tools.
type Registry struct {
	tools []Tool
}

func NewRegistry(tools ...Tool) *Registry {
	return &Registry{tools: tools}
}

func (r *Registry) All() []Tool {
	return r.tools
}

// Dispatch executes a named tool with the given input.
func (r *Registry) Dispatch(ctx context.Context, name string, input map[string]any) ToolResult {
	for _, t := range r.tools {
		if t.Name() == name {
			return t.Execute(ctx, input)
		}
	}
	return ToolResult{IsError: true, Content: "unknown tool: " + name}
}
