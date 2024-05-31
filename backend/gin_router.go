package main

import (
	"context"
	"encoding/json"
	"github.com/anda-ai/anda/agent"
	"github.com/anda-ai/anda/conf"
	"github.com/anda-ai/anda/entity"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
	"strings"
)

func router(cfg *conf.Config, r *gin.Engine) {

	router := NewRouter(cfg)

	chatRoute := r.Group("/")
	{
		chatRoute.GET("/chat", router.chat)
	}
}

type Router struct {
	chatAgent *agent.ChatAgent
	buffer    *websocket.Upgrader
}

func NewRouter(cfg *conf.Config) *Router {
	chatAgent := agent.NewChatAgent(cfg)
	return &Router{
		chatAgent: chatAgent,
		buffer: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (r *Router) chat(c *gin.Context) {
	conn, err := r.buffer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("upgrade err:%s", err)
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Warnf("close conn err:%s", err)
		}
	}(conn)

	ctx := context.Background()
	for {
		_, p, err := conn.ReadMessage()
		if err == nil {
			logger.Errorf("read msg err:%s", err)
			break
		}

		var query *entity.ChatQueryRequest

		if err := json.Unmarshal(p, query); err != nil {
			logger.Errorf("unmarshal err:%s", err)
			return
		}

		if strings.TrimSpace(query.Query) == "" {
			//TODO return others?
			return
		}

		resp, err := r.chatAgent.Searcher.Search(ctx, query.Query, 10)
		if err != nil {
			if err := writeErr(conn, err); err != nil {
				logger.Errorf("write search err:%s", err)
				return
			}
			return
		}

		if err := writeSearch(conn, resp); err != nil {
			logger.Errorf("write search err:%s", err)
			return
		}

		answer, err := r.chatAgent.Llm.ChatCompletion(ctx, makeMessages(ctx, resp, query))
		if err != nil {
			if err := writeErr(conn, err); err != nil {
				logger.Errorf("write err:%s", err)
				return
			}
			return
		}

		if err := writeAnswer(conn, answer); err != nil {
			logger.Errorf("write err:%s", err)
			return
		}

	}
}

func makeMessages(ctx context.Context, resp *entity.SearchResult, query *entity.ChatQueryRequest) []*entity.Message {
	msgs := make([]*entity.Message, 0, len(resp.Hits)+2)
	for _, hit := range resp.Hits {
		msgs = append(msgs, &entity.Message{
			Role:    "user",
			Content: hit.Snippet,
		})
	}

	msgs = append(msgs, &entity.Message{
		Role:    "user",
		Content: query.Query,
	})
	return msgs

}

func writeAnswer(conn *websocket.Conn, answer string) error {
	body, err := json.Marshal(&Response{
		Data:   answer,
		Status: 2,
	})
	if err != nil {
		return err

	}
	return conn.WriteMessage(websocket.TextMessage, body)
}

func writeSearch(conn *websocket.Conn, resp *entity.SearchResult) error {
	body, err := json.Marshal(&Response{
		Data:   resp,
		Status: 1,
	})
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, body)
}

type Response struct {
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
}

func writeErr(conn *websocket.Conn, e error) error {
	logger.Errorf("write err:%s", e)
	body, err := json.Marshal(&Response{
		Msg:    e.Error(),
		Status: 0,
	})
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, body)
}
