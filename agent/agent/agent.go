package agent

import (
	"context"
	"fmt"

	"github.com/awwwkshay/alpha-aide/llm"
)

const maxIter = 20

type Agent struct {
	provider llm.Provider
	tools    []llm.Tool
	history  []llm.Message
}

func New(provider llm.Provider, tools []llm.Tool) *Agent {
	return &Agent{provider: provider, tools: tools}
}

func NewWithHistory(provider llm.Provider, tools []llm.Tool, history []llm.Message) *Agent {
	return &Agent{provider: provider, tools: tools, history: history}
}

// OnChunk is called for each streaming chunk.
type OnChunk func(chunk llm.StreamChunk)

// Run processes a user message through the ReAct loop.
func (a *Agent) Run(ctx context.Context, userInput string, onChunk OnChunk) error {
	a.history = append(a.history, llm.Message{
		Role:    llm.RoleUser,
		Content: []llm.ContentBlock{{Type: "text", Text: userInput}},
	})

	for iter := 0; iter < maxIter; iter++ {
		ch, err := a.provider.Complete(ctx, SystemPrompt, a.history, a.tools)
		if err != nil {
			return fmt.Errorf("provider error: %w", err)
		}

		var assistantText string
		var toolCalls []*llm.ToolCall
		var streamErr error

		for chunk := range ch {
			switch chunk.Type {
			case "text_delta":
				assistantText += chunk.Text
				if onChunk != nil {
					onChunk(chunk)
				}
			case "tool_use_start":
				toolCalls = append(toolCalls, chunk.ToolCall)
				if onChunk != nil {
					onChunk(chunk)
				}
			case "tool_use_end":
				if onChunk != nil {
					onChunk(chunk)
				}
			case "error":
				streamErr = chunk.Err
			}
		}

		if streamErr != nil {
			return streamErr
		}

		// Build assistant message
		assistantMsg := llm.Message{Role: llm.RoleAssistant}
		if assistantText != "" {
			assistantMsg.Content = append(assistantMsg.Content, llm.ContentBlock{
				Type: "text",
				Text: assistantText,
			})
		}
		for _, tc := range toolCalls {
			assistantMsg.Content = append(assistantMsg.Content, llm.ContentBlock{
				Type:      "tool_use",
				ToolUseID: tc.ID,
				ToolName:  tc.Name,
				ToolInput: tc.Input,
			})
		}
		a.history = append(a.history, assistantMsg)

		if len(toolCalls) == 0 {
			return nil
		}

		// Execute tools and collect results
		resultMsg := llm.Message{Role: llm.RoleUser}
		for _, tc := range toolCalls {
			result := a.dispatchTool(ctx, tc)
			resultMsg.Content = append(resultMsg.Content, llm.ContentBlock{
				Type:      "tool_result",
				ToolUseID: tc.ID,
				ToolName:  tc.Name,
				Content:   result.Content,
				IsError:   result.IsError,
			})
			if onChunk != nil {
				onChunk(llm.StreamChunk{
					Type: "tool_result",
					Text: result.Content,
					ToolCall: &llm.ToolCall{
						ID:   tc.ID,
						Name: tc.Name,
					},
				})
			}
		}
		a.history = append(a.history, resultMsg)
	}

	return fmt.Errorf("max iterations (%d) reached", maxIter)
}

func (a *Agent) dispatchTool(ctx context.Context, tc *llm.ToolCall) llm.ToolResult {
	for _, t := range a.tools {
		if t.Name() == tc.Name {
			return t.Execute(ctx, tc.Input)
		}
	}
	return llm.ToolResult{IsError: true, Content: "unknown tool: " + tc.Name}
}

func (a *Agent) ClearHistory() {
	a.history = nil
}

func (a *Agent) History() []llm.Message {
	return a.history
}

func (a *Agent) ToolDefs() []llm.Tool {
	return a.tools
}
