package llms

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/entity"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseurl = "https://api.moonshot.cn/v1"

type MoonshotLLM struct {
	model       string
	apiKey      string
	Temperature float32
	client      *http.Client
}

type MoonshotMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MoonshotRequestBody struct {
	Model       string            `json:"model"`
	Messages    []MoonshotMessage `json:"messages"`
	Temperature float32           `json:"temperature"`
	Stream      bool              `json:"stream"`
}

func NewMoonshotLLM(cfg *conf.MoonshotCfg) LLM {

	logger.Infof("new moonshot api key: xxxxxxxxxx with model: %s temperature: %f timeout: %d s", cfg.Model, cfg.Temperature, cfg.TimeoutSec)

	key := os.Getenv("LLM_API_KEY")
	if key == "" {
		key = cfg.APIKey
	}

	return &MoonshotLLM{
		model:       cfg.Model,
		apiKey:      key,
		Temperature: cfg.Temperature,
		client: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSec) * time.Second,
		},
	}
}

func (o *MoonshotLLM) TokenCount(ctx context.Context, msg string) (int, error) {

	body := fmt.Sprintf(`{"model": "%s","messages": [{ "role": "user", "content": "%s" }]}`, o.model, msg)
	respBody, err := o.send(ctx, "/tokenizers/estimate-token-count", []byte(body))
	if err != nil {
		return 0, err
	}

	type Response struct {
		Code int `json:"code"`
		Data struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"data"`
		Scode  string `json:"scode"`
		Status bool   `json:"status"`
	}

	var result Response

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0, err
	}

	if !result.Status {
		return 0, fmt.Errorf("status is false")
	}

	return result.Data.TotalTokens, nil

}

func (o *MoonshotLLM) ChatCompletion(ctx context.Context, msgs []*entity.Message) (string, error) {

	req := MoonshotRequestBody{
		Model:       o.model,
		Messages:    make([]MoonshotMessage, 0),
		Temperature: o.Temperature,
	}

	for _, msg := range msgs {
		if msg.Content == "" {
			continue
		}

		req.Messages = append(req.Messages, MoonshotMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	if len(req.Messages) == 0 {
		return "", nil
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := o.send(ctx, "/chat/completions", body)
	if err != nil {
		return "", err
	}

	//{"id":"chatcmpl-2a6004b997b8408fb3670abed8ba0b93","object":"chat.completion","created":1716639433,"model":"moonshot-v1-8k","choices":[{"index":0,"message":{"role":"assistant","content":"你好，李雷！1+1等于2。"},"finish_reason":"stop"}],"usage":{"prompt_tokens":89,"completion_tokens":12,"total_tokens":101}}
	result := gjson.GetBytes(resp, "choices.0.message.content")
	if !result.Exists() {
		return "", fmt.Errorf("no completion found")
	}

	return result.String(), nil

}

func (o MoonshotLLM) ChatCompletionStream(ctx context.Context, msgs []*entity.Message) (*Stream, error) {
	req := MoonshotRequestBody{
		Model:       o.model,
		Messages:    make([]MoonshotMessage, 0),
		Temperature: o.Temperature,
		Stream:      true,
	}

	for _, msg := range msgs {
		if msg.Content == "" {
			continue
		}

		req.Messages = append(req.Messages, MoonshotMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no message")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	stream := NewStream()

	go func() {
		if err := o.streamSend(ctx, stream, "/chat/completions", body); err != nil {
			stream.error <- err
		}
	}()

	return stream, nil
}

func (o *MoonshotLLM) streamSend(ctx context.Context, stream *Stream, uri string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", baseurl+uri, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(resp.Body)

	br := bufio.NewReader(resp.Body)
	for {
		temp, _, e := br.ReadLine()
		if e != nil {
			if e == io.EOF {
				return nil
			} else {
				return err
			}
		}

		if len(temp) == 0 {
			continue
		}

		line := string(temp)

		if !strings.HasPrefix(line, "data:") {
			return fmt.Errorf("data: err line %s", line)
		}

		result := gjson.GetBytes(temp[5:], "choices.0.delta.content").String()

		if result == "" {
			continue
		}

		stream.inner <- result

	}
}

func (o *MoonshotLLM) send(ctx context.Context, uri string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", baseurl+uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s body: %s", resp.Status, string(respBody))
	}

	return respBody, nil
}
