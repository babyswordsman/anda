package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)
func router(r *gin.Engine){
	chatRoute := r.Group("/")
	{
		chatRoute.GET("/chat/stream",ChatStream)
	}
}


func ChatStream(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("panic %s \r\n %s", err, GetStack())
		}
	}()

	ctx.JSON(http.StatusOK,gin.H{
		"name":"anda ai",
	})
}