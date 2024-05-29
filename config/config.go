package config

import (
	"time"
)

type Config struct {
	Model         string
	FeedsFile     string
	OpenAIBaseUrl string
	FeedLimit     int
	LookBackDate  time.Time
	RankDropout   int
}

func DefaultConfig() Config {
	return Config{
		Model:         "llama3-70b-8192",
		OpenAIBaseUrl: "https://api.groq.com/openai/v1",
		FeedsFile:     "feeds",
		FeedLimit:     20,
		LookBackDate:  time.Now().AddDate(0, 0, -7),
		RankDropout:   2,
	}
}
