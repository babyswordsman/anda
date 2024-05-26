package calm

import (
	"context"
	"github.com/anda-ai/anda/models"
)

// Store is the interface for storing and searching messages
type Store interface {
	AddMessage(ctx context.Context, user string, message *models.Message) error
	Search(ctx context.Context, user, query string, pageSize int) ([]*models.Message, error)
	LastMessage(ctx context.Context, user string, size int) (*models.Message, error)
}
