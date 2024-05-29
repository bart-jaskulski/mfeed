package feed

import (
	"errors"
	"time"

	"github.com/mmcdole/gofeed"
)

var cutOffTime = time.Now().AddDate(0, 0, -7)

// FetchFeed fetches and parses an RSS feed from the given URL
func FetchFeed(url string) ([]FeedItem, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	if feed.UpdatedParsed != nil && !isNew(feed.UpdatedParsed) {
		return nil, errors.New("No new items")
	}

	var items []FeedItem
	for i, item := range feed.Items {
		if item.UpdatedParsed != nil && !isNew(item.UpdatedParsed) {
			continue
		}

		if item.PublishedParsed != nil && !isNew(item.PublishedParsed) {
			continue
		}

		var itemTime *time.Time
		now := time.Now()
		if item.PublishedParsed != nil {
			itemTime = item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			itemTime = item.UpdatedParsed
		} else {
			itemTime = &now
		}

		items = append(items, FeedItem{
			ID:          i,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source:      feed.Title,
			Updated:     itemTime,
		})
	}

	return items, nil
}

func isNew(item *time.Time) bool {
	if item == nil {
		return false
	}
	return item.After(cutOffTime)
}
