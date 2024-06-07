package calm

import (
	"context"
	"github.com/anda-ai/anda/entity"
)

// Store is the interface for storing and searching messages
type Store interface {
	AddMessage(ctx context.Context, user string, message *entity.Message) error
	Search(ctx context.Context, user, query string, pageSize int) ([]*entity.Message, error)
	LastMessage(ctx context.Context, user string, size int) (*entity.Message, error)
}
