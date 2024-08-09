package main

import (
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/chatroom"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/init"
	"likeminds-pandemonium/pubsub"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	initGin()
	redisClient := pubsub.InitRedisClient()

	router.Use(pubsub.ApiMiddleware(redisClient))
	// Application health check GET path
	router.GET("", init.HealthCheck)
	// ChatroomListen GET request
	router.GET(constant.ChatroomListen, chatroom.WsHandler())
	// PubSub publish / subscribe APIs
	router.POST(constant.RedisPublish, pubsub.Publish)

	// start server
	log.Fatal(router.Run(common.GoDotEnvVariable("Ws_SERVER_PORT")))
}

// initGin to initialise Gin network module
func initGin() {
	gin.SetMode(common.GoDotEnvVariable("GIN_MODE"))
	router = gin.Default()
}
