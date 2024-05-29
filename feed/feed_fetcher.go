package feed

import (
	"errors"
	"time"

	"github.com/mmcdole/gofeed"
)

// FetchFeed fetches and parses an RSS feed from the given URL
func FetchFeed(url string, lookBackTime time.Time) ([]FeedItem, error) {
	parser := gofeed.NewParser()
	feed, parseErr := parser.ParseURL(url)
	if parseErr != nil {
		return nil, parseErr
	}

	freshItems := filterNewItems(feed, lookBackTime)
	if len(freshItems) == 0 {
		return nil, errors.New("no new items found")
	}

  if len(freshItems) > 20 {
    freshItems = freshItems[:20]
  }

	return freshItems, nil
}

func filterNewItems(feed *gofeed.Feed, lookBackTime time.Time) []FeedItem {
	var freshItems []FeedItem
	for i, item := range feed.Items {
		itemTime := extractPublishTime(item)
		if itemTime == nil || !itemTime.After(lookBackTime) {
			continue
		}

		freshItems = append(freshItems, FeedItem{
			ID:          i,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source:      feed.Title,
			Updated:     itemTime,
		})
	}
	return freshItems
}

func extractPublishTime(item *gofeed.Item) *time.Time {
	if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	}
	now := time.Now()
	return &now
}
