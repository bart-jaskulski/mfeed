package config

import (
	"time"
)

const DefaultLookBackDays = -7 // Constant look-back days

type Config struct {
	LLMModel       string
	FeedsFilePath  string
	OpenAIEndpoint string
	MaxFeeds       int
	HistoricalDate time.Time
	MinimumRank    int
}

func NewConfig() Config {
	return Config{
		LLMModel:       "llama3-70b-8192",
		OpenAIEndpoint: "https://api.groq.com/openai/v1",
		FeedsFilePath:  "feeds",
		MaxFeeds:       20,
		HistoricalDate: time.Now().AddDate(0, 0, DefaultLookBackDays),
		MinimumRank:    3,
	}
}
