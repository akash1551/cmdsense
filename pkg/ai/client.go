package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cmdsense/pkg/config"
)

type Client interface {
	Explain(ctx context.Context, command string) (string, error)
}

func NewClient(cfg config.Config) (Client, error) {
	switch cfg.Provider {
	case "openai":
		if cfg.OpenAIKey == "" {
			return nil, fmt.Errorf("OpenAI API key is required. Run 'cmdsense config --set-openai-key <key>'")
		}
		return &OpenAIClient{APIKey: cfg.OpenAIKey, Model: cfg.Model}, nil
	case "ollama":
		return &OllamaClient{BaseURL: cfg.OllamaBaseURL, Model: cfg.Model}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

// --- OpenAI Client ---

type OpenAIClient struct {
	APIKey string
	Model  string
}

type openAIRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) Explain(ctx context.Context, command string) (string, error) {
	prompt := fmt.Sprintf("Explain this shell command concisely in Markdown. Warn if destructive. Command: `%s`", command)

	reqBody := openAIRequest{
		Model: c.Model,
		Messages: []message{
			{Role: "system", Content: "You are a helpful CLI assistant. Explain commands clearly using Markdown. Use 🔴 for dangerous commands."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}

// --- Ollama Client ---

type OllamaClient struct {
	BaseURL string
	Model   string
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (c *OllamaClient) Explain(ctx context.Context, command string) (string, error) {
	prompt := fmt.Sprintf("Explain this shell command concisely in Markdown. Warn if destructive. Command: `%s`", command)

	// Check if model exists (optional, could skip to generation)
	// For now, assume model is pulled.

	reqBody := ollamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/api/generate", strings.TrimSuffix(c.BaseURL, "/"))
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second} // Local LLMs can be slower
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Ollama connection error. Is it running? %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API error: status %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}
