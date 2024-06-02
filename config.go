package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const DefaultLookBackDays = -7 // Constant look-back days

type Config struct {
	LLMModel       string
	OpenAIEndpoint string
	FeedsFilePath  string
	MaxFeeds       int
	HistoricalDate time.Time
	MinimalScore   int
	FeedTitle      string
	FeedHref       string
	FeedAuthor     string
}

func defaultConfig() Config {
	return Config{
		LLMModel:       "llama3-70b-8192",
		OpenAIEndpoint: "https://api.groq.com/openai/v1",
		FeedsFilePath:  "feeds",
		MaxFeeds:       120,
		HistoricalDate: time.Now().AddDate(0, 0, DefaultLookBackDays),
		MinimalScore:   3,
		FeedTitle:      "Most interesting tech news",
		FeedHref:       "https://bart-jaskulski.github.io/mfeed/",
		FeedAuthor:     "Bart Jaskulski",
	}
}

func readFeedsFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file %s: %w", filename, err)
	}

	return lines, nil
}
