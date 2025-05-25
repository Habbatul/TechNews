package data

import "time"

type FeedItem struct {
	Title     string    `json:"title"`
	Link      string    `json:"link"`
	Published time.Time `json:"published"`
	Source    string    `json:"source"`
}
