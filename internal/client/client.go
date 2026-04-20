package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	anthropicAPIVersion = "2023-06-01"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents the API request
type Request struct {
	Model         string    `json:"model"`
	MaxTokens     int       `json:"max_tokens"`
	Messages      []Message `json:"messages"`
	Stream        bool      `json:"stream"`
	SystemPrompt  string    `json:"system,omitempty"`
}

// Response represents the API response chunk
type Response struct {
	Type     string `json:"type"`
	Index    int    `json:"index,omitempty"`
	Content  []struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"content,omitempty"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"delta,omitempty"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Client is the Anthropic API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{},
		model:      "claude-3-5-haiku-20241007",
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

// SendMessage sends a message and streams the response
func (c *Client) SendMessage(messages []Message, systemPrompt string) error {
	reqBody := Request{
		Model:        c.model,
		MaxTokens:    4096,
		Messages:     messages,
		Stream:       true,
		SystemPrompt: systemPrompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)

	fmt.Print("\n\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read response: %w", err)
		}

		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var response Response
		if err := json.Unmarshal([]byte(data), &response); err != nil {
			continue
		}

		if response.Error != nil {
			return fmt.Errorf("API error: %s", response.Error.Message)
		}

		if response.Type == "content_block_delta" {
			fmt.Print(response.Delta.Text)
		}
	}

	fmt.Println("\n")
	return nil
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