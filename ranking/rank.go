package ranking

import (
	"fmt"
	"sort"

	"github.com/bart-jaskulski/mfeed/config"
	"github.com/bart-jaskulski/mfeed/feed"
)

type Ranking struct {
	Items []feed.FeedItem
}

// CombineRankings merges multiple rankings into a single ranking
func CombineRankings(rankings ...Ranking) Ranking {
	var combined Ranking
	for _, r := range rankings {
		combined.Items = append(combined.Items, r.Items...)
	}
	return combined
}

// Rank ranks the items in the ranking
func (r *Ranking) Rank(cfg *config.Config) error {
	if len(r.Items) == 0 {
		return fmt.Errorf("no items to rank")
	}

	client := NewOpenAIClient(cfg)
	scoredItems, err := client.ScoreArticles(r.Items)
	if err != nil {
		return fmt.Errorf("failed to score articles: %w", err)
	}

	r.Items = scoredItems
	r.dropUninterestingItems(cfg)
	r.sortItems()
	r.Items = r.pickTopItems(cfg)

	return nil
}

func (r *Ranking) QuickRank(cfg *config.Config) error {
	if len(r.Items) == 0 {
		return fmt.Errorf("no items to rank")
	}

	r.dropUninterestingItems(cfg)
	r.sortItems()
	r.Items = r.pickTopItems(cfg)

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
		if item.Rank >= cfg.MinimumRank {
			filteredItems = append(filteredItems, item)
		}
	}
	r.Items = filteredItems
}

func (r *Ranking) pickTopItems(cfg *config.Config) []feed.FeedItem {
	topItems := r.Items
	if len(r.Items) > cfg.MaxFeeds {
		topItems = r.Items[:cfg.MaxFeeds]
	}
	return topItems
}
