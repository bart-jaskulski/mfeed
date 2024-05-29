package feed

import (
	"encoding/xml"
	"time"
)

// GenerateAtom generates an Atom feed with the top-ranked items
func GenerateAtom(items []FeedItem) (string, error) {
	feed := Feed{
		Xmlns: "http://www.w3.org/2005/Atom",
		Title: "Aggregated Feed",
		Link: struct {
			Href string `xml:"href,attr"`
		}{
			Href: "http://example.com",
		},
		Updated: time.Now().Format(time.RFC3339),
		ID:      "http://example.com",
		Author: AtomAuthor{
			Name: "Aggregated Feed Author",
		},
	}

	for _, item := range items {
		entry := AtomEntry{
			Title: item.Title,
			Link: struct {
				Href string `xml:"href,attr"`
			}{
				Href: item.Link,
			},
			ID:      item.Link,
			Updated: item.Updated.Format(time.RFC3339),
			Summary: item.Description,
			Author: AtomAuthor{
				Name: item.Source,
			},
			Rank: item.Rank,
		}
		feed.Entry = append(feed.Entry, entry)
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}
