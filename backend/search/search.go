package search

import (
	"context"
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/entity"
)

// Searcher is the interface for searching google bing or other search outer search engine ....
type Searcher interface {
	Search(ctx context.Context, query string, pageSize int) (*entity.SearchResult, error)
}

func NewSearcher(cfg *conf.SearchConfig) Searcher {
	if cfg.Serper != nil {
		return NewSerper(cfg.Serper)
	}
	return nil
}
