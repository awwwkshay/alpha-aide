package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/awwwkshay/alpha-aide/llm"
)

type WriteFileTool struct{}

func (WriteFileTool) Name() string        { return "write_file" }
func (WriteFileTool) Description() string { return "Write content to a file, creating parent directories as needed." }
func (WriteFileTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to write.",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "Content to write to the file.",
			},
		},
		"required": []string{"path", "content"},
	}
}

func (WriteFileTool) Execute(_ context.Context, input map[string]any) llm.ToolResult {
	path, _ := input["path"].(string)
	content, _ := input["content"].(string)
	if path == "" {
		return llm.ToolResult{IsError: true, Content: "path is required"}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("invalid path: %v", err)}
	}

	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("mkdir error: %v", err)}
	}

	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("write error: %v", err)}
	}

	return llm.ToolResult{Content: fmt.Sprintf("Written %d bytes to %s", len(content), path)}
}
