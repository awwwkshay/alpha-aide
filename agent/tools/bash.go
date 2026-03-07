package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
	"unicode/utf8"

	"github.com/awwwkshay/alpha-aide/llm"
)

const (
	maxOutputBytes = 100 * 1024 // 100KB per stream
	defaultTimeout = 30
	maxTimeout     = 300
)

type BashTool struct{}

func (BashTool) Name() string        { return "bash" }
func (BashTool) Description() string { return "Execute a bash command and return stdout, stderr, and exit code." }
func (BashTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The bash command to execute.",
			},
			"timeout_seconds": map[string]any{
				"type":        "integer",
				"description": "Timeout in seconds (default 30, max 300).",
			},
		},
		"required": []string{"command"},
	}
}

func (BashTool) Execute(ctx context.Context, input map[string]any) llm.ToolResult {
	command, _ := input["command"].(string)
	if command == "" {
		return llm.ToolResult{IsError: true, Content: "command is required"}
	}

	timeoutSec := defaultTimeout
	if v, ok := input["timeout_seconds"]; ok {
		switch t := v.(type) {
		case float64:
			timeoutSec = int(t)
		case int:
			timeoutSec = t
		}
	}
	if timeoutSec <= 0 {
		timeoutSec = defaultTimeout
	}
	if timeoutSec > maxTimeout {
		timeoutSec = maxTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else if ctx.Err() != nil {
			return llm.ToolResult{Content: fmt.Sprintf("[exit code: -1]\n--- stdout ---\n%s\n--- stderr ---\ntimeout after %ds",
				sanitize(stdout.Bytes()), timeoutSec)}
		}
	}

	return llm.ToolResult{
		Content: fmt.Sprintf("[exit code: %d]\n--- stdout ---\n%s\n--- stderr ---\n%s",
			exitCode,
			sanitize(stdout.Bytes()),
			sanitize(stderr.Bytes()),
		),
	}
}

func sanitize(b []byte) string {
	if len(b) > maxOutputBytes {
		b = b[:maxOutputBytes]
	}
	if !utf8.Valid(b) {
		b = bytes.ToValidUTF8(b, []byte("?"))
	}
	return string(b)
}
