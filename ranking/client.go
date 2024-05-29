package ranking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/bart-jaskulski/feed/feed"
	openai "github.com/sashabaranov/go-openai"
)

// OpenAIClient wraps the OpenAI API client
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient() *OpenAIClient {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatalf("GROQ_API_KEY environment variable not set")
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"
	return &OpenAIClient{openai.NewClientWithConfig(config)}
}

type Scoring struct {
	Articles []ItemScore `json:"articles"`
}

type ItemScore struct {
	ID    int     `json:"id"`
	Score float64 `json:"score"`
}

// ScoreArticles sends item titles to OpenAI API for scoring
func (c *OpenAIClient) ScoreArticles(items []feed.FeedItem) ([]feed.FeedItem, error) {
	ctx := context.Background()

	rankPrompt, err := preparePrompt("prompts/rank.tmpl", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare prompt: %w", err)
	}

	data := struct {
		Items []feed.FeedItem
	}{
		Items: items,
	}
	itemsPrompt, err := preparePrompt("prompts/titles.tmpl", data)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare prompt: %w", err)
	}

	req := openai.ChatCompletionRequest{
		Model: "llama3-70b-8192",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: rankPrompt,
			},
			{
				Role:    "user",
				Content: itemsPrompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.0,
		MaxTokens:   1000,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	var scores Scoring

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &scores)
	if err != nil {
		return nil, err
	}

	for i := range items {
		for _, score := range scores.Articles {
			if items[i].ID == score.ID {
				items[i].Rank = int(score.Score)
				break
			}
		}
	}

	return items, nil
}

func preparePrompt(promptFile string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt template: %w", err)
	}
	tmplBuffer := bytes.Buffer{}
	err = tmpl.Execute(&tmplBuffer, data)
	if err != nil {
		return "", err
	}
	return tmplBuffer.String(), nil
}
