package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/awwwkshay/alpha-aide/llm"
)

type ReadFileTool struct{}

func (ReadFileTool) Name() string        { return "read_file" }
func (ReadFileTool) Description() string { return "Read the contents of a file." }
func (ReadFileTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to read.",
			},
		},
		"required": []string{"path"},
	}
}

func (ReadFileTool) Execute(_ context.Context, input map[string]any) llm.ToolResult {
	path, _ := input["path"].(string)
	if path == "" {
		return llm.ToolResult{IsError: true, Content: "path is required"}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("invalid path: %v", err)}
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("read error: %v", err)}
	}

	// Binary file heuristic
	if bytes.Contains(data, []byte{0}) {
		return llm.ToolResult{IsError: true, Content: "file appears to be binary"}
	}

	return llm.ToolResult{Content: fmt.Sprintf("--- %s ---\n%s", path, string(data))}
}
