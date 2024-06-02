package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/bart-jaskulski/mfeed/config"
	"github.com/bart-jaskulski/mfeed/feed"
	"github.com/bart-jaskulski/mfeed/fileutil"
	"github.com/bart-jaskulski/mfeed/ranking"
	readability "github.com/go-shiori/go-readability"
)

func main() {
	cfg := config.NewConfig()

	urls, readErr := fileutil.ReadLines(cfg.FeedsFilePath)
	if readErr != nil {
		log.Fatalf("error reading feed URLs: %v", readErr)
	}

	if len(urls) == 0 {
		log.Fatalf("no feed URLs found")
	}

	var wg sync.WaitGroup
	rankingsChan := make(chan ranking.Ranking, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(feedURL string) {
			defer wg.Done()
			rank, err := processFeed(feedURL, &cfg)
			if err != nil {
        log.Printf("error generating feed: %v", err)
				return
			}
			if len(rank.Items) > 0 {
				rankingsChan <- rank
			}
		}(url)
	}

	wg.Wait()
	close(rankingsChan)

	var allRankings []ranking.Ranking
	for rank := range rankingsChan {
		allRankings = append(allRankings, rank)
	}

	if len(allRankings) == 0 {
		log.Fatalf("no valid rankings found")
	}

	metaRanking := ranking.CombineRankings(allRankings...)
	if processErr := metaRanking.QuickRank(&cfg); processErr != nil {
		log.Fatalf("error during ranking process: %v", processErr)
	}

	if len(metaRanking.Items) == 0 {
		log.Fatalf("no new items after ranking")
	}

	for i, feedItem := range metaRanking.Items {
		if feedItem.Content != "" {
			// content already available
			continue
		}

		feedItem.Content = createContent(feedItem, &cfg)

		if feedItem.Rank == cfg.MinimumRank {
			feedItem.Content = createSummary(feedItem, &cfg)
		}

		metaRanking.Items[i].Content = feedItem.Content
	}

	// Output the combined and ranked RSS feed
	feed, genErr := feed.GenerateFeed(metaRanking.Items)
	if genErr != nil {
		log.Fatalf("error generating feed: %v", genErr)
	}

	os.Stdout.Write([]byte(feed))
}

func createContent(feedItem feed.FeedItem, cfg *config.Config) string {
	art, err := readability.FromURL(feedItem.Link, 30*time.Second)
	if err != nil {
		log.Printf("error fetching article: %v", err)
		return ""
	}

	return art.Content
}

func createSummary(feedItem feed.FeedItem, cfg *config.Config) string {
	ai := ranking.NewOpenAIClient(cfg)
	summary, err := ai.SummarizeContent(feedItem.Content)
	if err != nil {
		log.Printf("error summarizing content: %v", err)
		return ""
	}
	return summary
}

func processFeed(feedURL string, cfg *config.Config) (ranking.Ranking, error) {
	items, fetchErr := feed.FetchFeed(feedURL, cfg.HistoricalDate)
	if fetchErr != nil {
		return ranking.Ranking{}, fetchErr
	}

	if len(items) == 0 {
		return ranking.Ranking{}, nil
	}

	feedRanking := ranking.Ranking{Items: items}
	if rankErr := feedRanking.Rank(cfg); rankErr != nil {
		return ranking.Ranking{}, rankErr
	}

	return feedRanking, nil
}
