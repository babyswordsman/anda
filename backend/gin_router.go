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

	logger.Infof("chat recive conn on %s", c.Request.Host)

	conn, err := r.buffer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("upgrade err:%s", err)
		return
	}

	defer func(conn *websocket.Conn) {
		if conn == nil {
			return
		}
		err := conn.Close()
		if err != nil {
			logger.Warnf("close conn err:%s", err)
		}
	}(conn)

	ctx := context.Background()
	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			logger.Errorf("read msg err:%s", err)
			break
		}

		bs, err := json.Marshal(&entity.ChatQueryRequest{
			Query: "中国的国土面积有多少大",
		})

		logger.Infof("read message query:[%s]", string(bs))

		logger.Infof("read message query:[%s]", string(p))

		var query = &entity.ChatQueryRequest{}
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

		stream, err := r.chatAgent.Llm.ChatCompletionStream(ctx, makeMessages(ctx, resp, query))
		if err != nil {
			logger.Errorf("read answer has err:%s", err)
			if err := writeErr(conn, err); err != nil {
				logger.Errorf("write err:%s", err)
				return
			}
			return
		}

		for {
			answer, err := stream.Next()

			if answer == "" && err == nil {
				break
			}

			if err != nil {
				logger.Errorf("read answer has err:%s", err)
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

		if err := writeAnswer(conn, "<END_EOF>"); err != nil {
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
	Data   interface{} `json:"data"`   // if status is 0, it write error message
	Msg    string      `json:"msg"`    // if status is 1, it's search result, if status is 2, it's answer
	Status int         `json:"status"` // 0 err , 1 search , 2 answer
}

func writeErr(conn *websocket.Conn, e error) error {
	body, err := json.Marshal(&Response{
		Msg:    e.Error(),
		Status: 0,
	})
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, body)
}
