package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	api   *openai.Client
	model string
}

func NewClient(apiKey string) *Client {
	// Default to GPT-4o
	return &Client{
		api:   openai.NewClient(apiKey),
		model: "gpt-4o", // Using string literal to be safe
	}
}

// GenerateJSON executes a chat completion ensuring JSON output
func (c *Client) GenerateJSON(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	}

	resp, err := c.api.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("openai completion error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from openai")
	}

	return resp.Choices[0].Message.Content, nil
}
