package redisPandemonium

import (
	"context"
	"github.com/gin-gonic/gin"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"log"
)

var ctx = context.Background()

func Publish(c *gin.Context) {
	Redis(c, constant.POSTMethod)
}
func Redis(c *gin.Context, method int) {
	switch method {
	case constant.POSTMethod:
		topic := c.Param("topic")
		if topic == "" || topic == "null" {
			api.GeneralBadRequestError(c, "topic is required")
			return
		}
		rawData, _ := c.GetRawData()
		redisClient := GetRedisClientFromContext(c)
		// publish rawData to redisPandemonium channel:<chatroomID>
		if err := redisClient.Publish(ctx, topic, rawData).Err(); err != nil {
			api.GeneralAPIError(c, err.Error())
			log.Println("Publish error:", err)
			return
		}
		api.GenerateResponse(c, nil)
	}
}
