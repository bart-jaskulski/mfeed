package main

import (
	"log"
	"os"
	"sync"

	"github.com/bart-jaskulski/feed/feed"
	"github.com/bart-jaskulski/feed/fileutil"
	"github.com/bart-jaskulski/feed/ranking"
	"github.com/bart-jaskulski/feed/config"
)

func main() {
  config := config.DefaultConfig()

	feedURLs, err := fileutil.ReadLines(config.FeedsFile)
	if err != nil {
		log.Fatalf("Error reading feed URLs from file: %v", err)
	}

	if len(feedURLs) == 0 {
		log.Fatalf("No feed URLs found in file")
	}

	var wg sync.WaitGroup
	rankingsChan := make(chan ranking.Ranking, len(feedURLs))

	for _, url := range feedURLs {
		wg.Add(1)
		go fetchAndRankFeed(url, &wg, rankingsChan, &config)
	}

	wg.Wait()
	close(rankingsChan)

  metaRanking := ranking.CombineRankings(rankingsChan...)
  metaRanking.Rank(&config)

	if len(metaRanking.Items) == 0 {
		log.Fatalf("No new items")
	}

	// Generate the RSS feed file
	rssContent, err := feed.GenerateAtom(metaRanking.Items)
	if err != nil {
		log.Fatalf("Error generating RSS: %v", err)
	}

	os.Stdout.Write([]byte(rssContent))
}

func fetchAndRankFeed(url string, wg *sync.WaitGroup, rankingsChan chan<- ranking.Ranking, config *config.Config) {
	defer wg.Done()

	items, err := feed.FetchFeed(url)
	if err != nil {
		log.Printf("Error fetching feed %s: %v", url, err)
		return
	}

	if len(items) == 0 {
		return
	}

  ranking := ranking.Ranking{Items: items}
  err = ranking.Rank(config)
	if err != nil {
		log.Printf("Error ranking items for feed %s: %v", url, err)
		return
	}

	rankingsChan <- ranking
}
