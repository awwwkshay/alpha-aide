package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/awwwkshay/alpha-aide/llm"
)

type EditFileTool struct{}

func (EditFileTool) Name() string        { return "edit_file" }
func (EditFileTool) Description() string { return "Replace a unique string in a file with new content." }
func (EditFileTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to edit.",
			},
			"old_str": map[string]any{
				"type":        "string",
				"description": "The exact string to replace. Must appear exactly once in the file.",
			},
			"new_str": map[string]any{
				"type":        "string",
				"description": "The string to replace old_str with.",
			},
		},
		"required": []string{"path", "old_str", "new_str"},
	}
}

func (EditFileTool) Execute(_ context.Context, input map[string]any) llm.ToolResult {
	path, _ := input["path"].(string)
	oldStr, _ := input["old_str"].(string)
	newStr, _ := input["new_str"].(string)

	if path == "" {
		return llm.ToolResult{IsError: true, Content: "path is required"}
	}
	if oldStr == "" {
		return llm.ToolResult{IsError: true, Content: "old_str is required"}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("invalid path: %v", err)}
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("read error: %v", err)}
	}

	content := string(data)
	count := strings.Count(content, oldStr)
	switch count {
	case 0:
		return llm.ToolResult{IsError: true, Content: "String not found. No changes made."}
	case 1:
		// proceed
	default:
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("Found %d times (must be unique). Add more context.", count)}
	}

	newContent := strings.Replace(content, oldStr, newStr, 1)
	if err := os.WriteFile(absPath, []byte(newContent), 0644); err != nil {
		return llm.ToolResult{IsError: true, Content: fmt.Sprintf("write error: %v", err)}
	}

	return llm.ToolResult{Content: fmt.Sprintf("Edited %s successfully.", path)}
}
