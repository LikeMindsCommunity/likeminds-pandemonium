package main

import (
	"github.com/gin-gonic/gin"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/pubsub"
	"log"
)

var (
	router *gin.Engine
)

func main() {
	initGin()
	redisClient := pubsub.InitRedisClient()

	router.Use(pubsub.ApiMiddleware(redisClient))
	// Subscribe
	router.GET(constant.RedisSubscribe, pubsub.Subscribe())
	// Publish
	router.POST(constant.RedisPublish, pubsub.Publish)

	// start server
	log.Fatal(router.Run(common.GoDotEnvVariable("Ws_SERVER_PORT")))
}

// initGin to initialise Gin network module
func initGin() {
	gin.SetMode(common.GoDotEnvVariable("GIN_MODE"))
	router = gin.Default()
}
