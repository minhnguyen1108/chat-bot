package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents the API request (OpenAI compatible)
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Response represents the API response
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Client is the AI Chat API client
type Client struct {
	apiKey string
	model  string
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  "anthropic/claude-sonnet-4.5",
	}
}

// SetModel sets the model to use
func (c *Client) SetModel(model string) {
	c.model = model
}

// GetModel returns the current model
func (c *Client) GetModel() string {
	return c.model
}

// SendMessage sends a message and returns the response
func (c *Client) SendMessage(messages []Message, systemPrompt string) (string, error) {
	// Add system prompt as first message if provided
	var allMessages []Message
	if systemPrompt != "" {
		allMessages = append(allMessages, Message{Role: "system", Content: systemPrompt})
	}
	allMessages = append(allMessages, messages...)

	reqBody := Request{
		Model:    c.model,
		Messages: allMessages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://aishop24h.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to decode response: %s", string(body))
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", nil
}

// GetAvailableModels returns available models
func GetAvailableModels() []string {
	return []string{
		"anthropic/claude-sonnet-4.5",
		"anthropic/claude-haiku-4",
		"anthropic/claude-opus-4",
	}
}
