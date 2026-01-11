package model

import "time"

type Item struct {
	Title      string    `json:"title"`
	Categories []string  `json:"categories"`
	Link       string    `json:"link"`
	Date       time.Time `json:"date"`
	Summary    string    `json:"summary"`
	SourceName string    `json:"source_name"`
}

type Source struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	FeedURL   string    `json:"feed_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Article struct {
	ID          int64
	SourceID    int64
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
	PostedAt    time.Time
	CreatedAt   time.Time
}
