package llms

import (
	"context"
	"github.com/anda-ai/anda/models"
)

type LLM interface {
	TokenCount(ctx context.Context, msg string) (int, error)

	ChatCompletion(ctx context.Context, request []*models.Message) (string, error)

	ChatCompletionStream(ctx context.Context, request []*models.Message) (*Stream, error)
}

type Stream struct {
	inner chan string
}
