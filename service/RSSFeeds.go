package service

import (
	"TechNews/data"
	"github.com/mmcdole/gofeed"
	"log"
	"sort"
	"time"
)

var feedSources = map[string]string{
	"hackernews":         "https://hnrss.org/frontpage",
	"codrops (dev.to)":   "https://dev.to/feed",
	"dev-community":      "https://tympanus.net/codrops/feed/",
	"stackoverflow-blog": "https://stackoverflow.blog/feed/",
	"verge":              "https://www.theverge.com/rss/index.xml",
	"techcrunch":         "https://techcrunch.com/feed/",
	"youtube-fireship":   "https://www.youtube.com/feeds/videos.xml?channel_id=UCsBjURrPoezykLs9EqgamOA",
	"github-trend":       "https://github-rss.alexi.sh/feeds/daily/all.xml",
	"youtube-codehead":   "https://www.youtube.com/feeds/videos.xml?channel_id=UCl6j6qHPHocyF9HSlsdQnqw",
	"TLDR-Tech":          "https://tldr.tech/api/rss/tech",
}

func fetchLatest(feedURL string, count int) []data.FeedItem {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(feedURL)
	if err != nil {
		log.Printf("Gagal fetch %s: %v", feedURL, err)
		return nil
	}

	var items []data.FeedItem
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	for _, i := range feed.Items {
		published := time.Now()
		if i.PublishedParsed != nil {
			published = *i.PublishedParsed
		}

		//7 hari terakhir
		if published.After(oneWeekAgo) {
			items = append(items, data.FeedItem{
				Title:     i.Title,
				Link:      i.Link,
				Published: published,
				Source:    feed.Title,
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Published.After(items[j].Published)
	})

	if len(items) > count {
		return items[:count]
	}
	return items
}

func GetNews() map[string][]data.FeedItem {
	result := make(map[string][]data.FeedItem)
	for key, url := range feedSources {
		result[key] = fetchLatest(url, 7)
	}

	return result
}
