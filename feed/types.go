package feed

import (
	"encoding/xml"
	"time"
)

type FeedItem struct {
	ID      int
	Title   string
	Link    string
	Content string
	Source  string
	Updated *time.Time
	Rank    int
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Updated string      `xml:"updated"`
	ID      string      `xml:"id"`
	Author  AtomAuthor  `xml:"author"`
	Entry   []AtomEntry `xml:"entry"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}

type AtomEntry struct {
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Content AtomContent `xml:"content,omitempty"`
	Author  AtomAuthor  `xml:"author"`
}

type AtomContent struct {
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}
