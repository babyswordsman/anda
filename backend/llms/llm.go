package llms

import (
	"context"
	"github.com/anda-ai/anda/entity"
)

type LLM interface {
	TokenCount(ctx context.Context, msg string) (int, error)

	ChatCompletion(ctx context.Context, request []*entity.Message) (string, error)

	ChatCompletionStream(ctx context.Context, request []*entity.Message) (*Stream, error)
}

type Stream struct {
	inner chan string
}
