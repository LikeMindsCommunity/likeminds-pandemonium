package main

import (
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/middleware"
	"likeminds-pandemonium/pubsub"
	"likeminds-pandemonium/web"
	"likeminds-pandemonium/ws"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

const (
	// AppVersion | current application version
	AppVersion string = "0.1.1"
)

func main() {
	initGin()
	redisClient := pubsub.InitRedisClient()
	wsServerParent := ws.NewWsServerParent()
	router.Use(middleware.ApiMiddleware(redisClient, wsServerParent))

	router.GET("", web.Home)
	router.GET(constant.RedisSubscribe, pubsub.Subscribe())
	// Publish
	router.POST(constant.RedisPublish, pubsub.Publish)
	// Delivery Report
	router.GET(constant.DeliveredDR, pubsub.DeliveryReportHandler)

	// start server
	log.Fatal(router.Run(common.GoDotEnvVariable("Ws_SERVER_PORT")))
}

// initGin to initialise Gin network module
func initGin() {
	gin.SetMode(common.GoDotEnvVariable("GIN_MODE"))
	router = gin.Default()
}
