package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const anthropicAPIVersion = "2023-06-01"

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents the API request
type Request struct {
	Model        string    `json:"model"`
	MaxTokens    int       `json:"max_tokens"`
	Messages     []Message `json:"messages"`
	Stream       bool      `json:"stream"`
	SystemPrompt string    `json:"system,omitempty"`
}

// Client is the Anthropic API client
type Client struct {
	apiKey string
	model  string
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  "claude-3-5-haiku-20241007",
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
	reqBody := Request{
		Model:        c.model,
		MaxTokens:    4096,
		Messages:     messages,
		Stream:       false,
		SystemPrompt: systemPrompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Content) > 0 {
		return response.Content[0].Text, nil
	}

	return "", nil
}

// GetAvailableModels returns available models
func GetAvailableModels() []string {
	return []string{
		"claude-3-5-haiku-20241007",
		"claude-3-5-sonnet-20241022",
		"claude-3-opus-20240229",
		"claude-3-haiku-20240307",
	}
}
