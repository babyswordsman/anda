package llms

import (
	"context"
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/entity"
	logger "github.com/sirupsen/logrus"
)

type LLM interface {
	TokenCount(ctx context.Context, msg string) (int, error)

	ChatCompletion(ctx context.Context, request []*entity.Message) (string, error)

	ChatCompletionStream(ctx context.Context, request []*entity.Message) (*Stream, error)
}

type Stream struct {
	inner chan string
}

func NewLLM(cfg *conf.LLMConfig) LLM {
	if cfg.Moonshot != nil {
		return NewMoonshotLLM(cfg.Moonshot)
	}

	logger.Errorf("not instance llm client cfg:%v", cfg)

	return nil

}
