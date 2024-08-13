package main

import (
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/pubsub"
	"likeminds-pandemonium/web"
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
	
	router.GET("", web.Home)
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
