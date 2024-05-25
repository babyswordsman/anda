package models

type SearchResult struct {
	hits            []*Hit
	RelatedSearches []string
}

type Hit struct {
	Title   string
	Link    string
	snippet string
}
