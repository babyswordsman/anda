package llms

import (
	"context"
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/entity"
	"os"
	"testing"
)

func TestMoonshotLLM_TokenCount(t *testing.T) {

	tests := []struct {
		msg     string
		want    int
		wantErr bool
	}{
		{
			msg:     "Hello, how are you?",
			want:    13,
			wantErr: false,
		},
		{
			msg:     "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。",
			want:    72,
			wantErr: false,
		},
	}

	llm := NewMoonshotLLM(&conf.MoonshotCfg{
		APIKey:     os.Getenv("LLM_API_KEY"),
		TimeoutSec: 10,
		Model:      "moonshot-v1-8k",
	})

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			ctx := context.Background()
			got, err := llm.TokenCount(ctx, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TokenCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoonshotLLM_ChatCompletion(t *testing.T) {

	llm := NewMoonshotLLM(&conf.MoonshotCfg{
		APIKey:      os.Getenv("LLM_API_KEY"),
		TimeoutSec:  10,
		Model:       "moonshot-v1-8k",
		Temperature: 0.3,
	})

	msgs := []*entity.Message{
		{
			Role:    "system",
			Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。",
		},
		{
			Role:    "user",
			Content: "你好，我叫李雷，1+1等于多少？",
		},
	}

	ctx := context.Background()
	got, err := llm.ChatCompletion(ctx, msgs)
	if err != nil {
		t.Fatal("ChatCompletion() error ", err)

	}
	println(got)
}
