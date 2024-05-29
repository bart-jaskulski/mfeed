package ranking

import (
	"fmt"
	"os"
	"sort"

	"github.com/bart-jaskulski/feed/feed"
	"github.com/bart-jaskulski/feed/config"
)

const RankDropout = 2

type Ranking struct {
  Items []feed.FeedItem
  Config *config.Config
}

func RankItems(items []feed.FeedItem) ([]feed.FeedItem, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	client := NewOpenAIClient()
	scoredItems, err := client.ScoreArticles(items)
	if err != nil {
		return nil, err
	}

	return scoredItems, nil
}

func CombineRankings(ranking ...Ranking) Ranking {
  var combined Ranking
  for _, r := range ranking {
    combined.Items = append(combined.Items, r.Items...)
  }
  return combined
}

func (r *Ranking) Rank(cfg *config.Config) error {
	client := NewOpenAIClient()
	scoredItems, err := client.ScoreArticles(r.Items)
  if err != nil {
    return err
  }

  r.Items = scoredItems
  r.dropUninterestingItems(cfg)
  r.sortItems()
  r.pickTopItems(cfg)

  return nil
}

func (r *Ranking) sortItems() {
  sort.Slice(r.Items, func(i, j int) bool {
    return r.Items[i].Rank > r.Items[j].Rank
  })
}

func (r *Ranking) dropUninterestingItems(cfg *config.Config) {
  var filteredItems []feed.FeedItem
  for _, item := range r.Items {
    if item.Rank > cfg.RankDropout {
      r.Items = append(filteredItems, item)
    }
  }
}

func (r *Ranking) pickTopItems(cfg *config.Config) []feed.FeedItem {
  topItems := r.Items
  if len(r.Items) > cfg.FeedLimit {
    topItems = r.Items[:r.Config.FeedLimit]
  }
  return topItems
}

func SortItems(items []feed.FeedItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Rank > items[j].Rank
	})
}

func DropUninterestingItems(items []feed.FeedItem) []feed.FeedItem {
	var filteredItems []feed.FeedItem
	for _, item := range items {
		if item.Rank > RankDropout {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func PickTopItems(items []feed.FeedItem) []feed.FeedItem {
	topItems := items
	if len(items) > 20 {
		topItems = items[:20]
	}

	return topItems
}
