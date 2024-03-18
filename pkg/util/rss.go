package util

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

// ParseFeed parses an RSS feed and returns the items
func ParseFeed(url string) ([]*gofeed.Item, error) {
	fp := gofeed.NewParser()

	// Parse the URL and ensure it's an RSS feed
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSS feed: %v", err)
	}

	return feed.Items, nil
}

// GetItemsBetweenDates returns the items that were published between the start and end dates
func GetItemsBetweenDates(items []*gofeed.Item, startDate, endDate time.Time) []*gofeed.Item {
	var filteredItems []*gofeed.Item
	for _, item := range items {
		if item.PublishedParsed.After(startDate) && item.PublishedParsed.Before(endDate) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}
