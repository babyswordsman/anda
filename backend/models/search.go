package models

import "time"

type SearchResult struct {
	Hits            []*Hit   `json:"hits"`
	RelatedSearches []string `json:"related_searches"`
}

type Hit struct {
	Title   string    `json:"title"`
	Link    string    `json:"link"`
	Snippet string    `json:"snippet"`
	Date    time.Time `json:"date"`
}
