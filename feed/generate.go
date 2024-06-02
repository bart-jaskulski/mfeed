package feed

import (
	"encoding/xml"
	"errors"
	"time"
)

// GenerateFeed generates an Atom feed with the top-ranked items
func GenerateFeed(items []FeedItem) (string, error) {
	if len(items) == 0 {
		return "", errors.New("no items to process")
	}

	feed := createAtomFeed(items)
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}

func createAtomFeed(items []FeedItem) Feed {
	feed := Feed{
		Xmlns: "http://www.w3.org/2005/Atom",
		Title: "Aggregation of Most Interesting Feeds",
		Link: struct {
			Href string `xml:"href,attr"`
		}{
			Href: "http://example.com",
		},
		Updated: time.Now().Format(time.RFC3339),
		ID:      "http://example.com",
		Author: AtomAuthor{
			Name: "Author",
		},
	}

	for _, item := range items {
		entry := AtomEntry{
			Title: fmt.Sprintf("[%d] %s", item.Rank, item.Title),
			Link: struct {
				Href string `xml:"href,attr"`
			}{
				Href: item.Link,
			},
			ID:      item.Link,
			Updated: item.Updated.Format(time.RFC3339),
			Content: AtomContent{
				Content: item.Content,
				Type:    "html",
			},
			Author: AtomAuthor{
				Name: item.Source,
			},
		}
		feed.Entry = append(feed.Entry, entry)
	}
	return feed
}
