package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/entity"
	logger "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

// Serper GET Google Search Results
// https://serper.dev
type Serper struct {
	cfg    *conf.SerperCfg
	client *http.Client
}

func NewSerper(cfg *conf.SerperCfg) Searcher {

	logger.Infof("new serper api key: xxxxxxxxxx with timeout: %d s", cfg.TimeoutSec)

	key := os.Getenv("SERPER_API_KEY")
	if key == "" {
		key = cfg.APIKey
	}

	logger.Infof("serper api key: %s with timeout: %d s", key, cfg.TimeoutSec)

	return &Serper{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSec) * time.Second,
		},
	}
}

func (s *Serper) Search(ctx context.Context, query string, pageSize int) (*entity.SearchResult, error) {
	bodyParam := make(map[string]interface{})
	bodyParam["q"] = query
	bodyParam["location"] = "China"
	bodyParam["gl"] = "cn"
	bodyParam["hl"] = "zh-cn"
	bodyParam["num"] = pageSize

	jsonData, err := json.Marshal(bodyParam)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://google.serper.dev/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", "fedf600ef33760e9c8bafb6cad43d5d77251145f")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s body: %s", resp.Status, string(respBody))
	}

	serperResult := &SerperSearchResult{}
	if err := json.Unmarshal(respBody, &serperResult); err != nil {
		return nil, err
	}

	result := &entity.SearchResult{}
	for _, item := range serperResult.RelatedSearches {
		result.RelatedSearches = append(result.RelatedSearches, item.Query)

	}

	for _, item := range serperResult.Organic {
		const layout = "2006年1月2日"
		parse, err := time.Parse(layout, item.Date)
		if err != nil {
			logger.Warnf("value:%s parse date error: %s", item.Date, err.Error())
		}
		result.Hits = append(result.Hits, &entity.Hit{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
			Date:    parse,
		})

	}
	return result, nil
}

type SerperSearchParameters struct {
	Q        string `json:"q"`
	Gl       string `json:"gl"`
	Hl       string `json:"hl"`
	Type     string `json:"type"`
	Num      int    `json:"num"`
	Location string `json:"location"`
	Engine   string `json:"engine"`
}

type SerperKnowledgeGraph struct {
	Title             string `json:"title"`
	Type              string `json:"type"`
	ImageUrl          string `json:"imageUrl"`
	Description       string `json:"description"`
	DescriptionSource string `json:"descriptionSource"`
	DescriptionLink   string `json:"descriptionLink"`
}

type SerperOrganic struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Snippet  string `json:"snippet"`
	Date     string `json:"date"`
	Position int    `json:"position"`
}

type SerperRelatedSearches struct {
	Query string `json:"query"`
}

type SerperSearchResult struct {
	SearchParameters SerperSearchParameters  `json:"searchParameters"`
	KnowledgeGraph   SerperKnowledgeGraph    `json:"knowledgeGraph"`
	Organic          []SerperOrganic         `json:"organic"`
	RelatedSearches  []SerperRelatedSearches `json:"relatedSearches"`
}
