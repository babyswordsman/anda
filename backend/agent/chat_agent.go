package agent

import (
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/llms"
	"github.com/anda-ai/anda/search"
)

type ChatAgent struct {
	Llm      llms.LLM
	Searcher search.Searcher
}

func NewChatAgent(conf *conf.Config) *ChatAgent {
	llm := llms.NewLLM(conf.LLMConfig)
	searcher := search.NewSearcher(conf.SearchConfig)
	return &ChatAgent{
		Llm:      llm,
		Searcher: searcher,
	}
}
