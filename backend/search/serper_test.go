package search

import (
	"context"
	"encoding/json"
	"github.com/anda-ai/anda/conf"
	"testing"
)

func TestSerper_Search(t *testing.T) {

	s := NewSerper(&conf.SerperCfg{
		APIKey:     "your api key here",
		TimeoutSec: 30,
	})

	got, err := s.Search(context.Background(), "中国的国土面积是多大", 2)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
		return
	}

	marshal, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	println(string(marshal))
}
