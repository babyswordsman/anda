package search

import (
	"context"
	"github.com/anda-ai/anda/models"
)

type Searcher interface {
	Search(ctx context.Context, sessionID, query string, pageSize int) ([]*models.SearchResult, error)
}
