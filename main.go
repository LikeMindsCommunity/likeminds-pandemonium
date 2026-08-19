package main

import (
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/constant"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/common"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/database"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/middleware"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/pubsub"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/web"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/ws"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

const (
	// AppVersion | current application version
	AppVersion string = "0.4.2"
)

func main() {
	initGin()
	redisClient := pubsub.InitRedisClient()
	database.Postgres = database.ConnectPostgres()
	wsServerParent := ws.NewWsServerParent()
	router.Use(middleware.ApiMiddleware(redisClient, wsServerParent))

	router.GET("", web.Home)
	router.GET(constant.RedisSubscribe, pubsub.Subscribe())
	// Publish
	router.POST(constant.RedisPublish, pubsub.Publish)
	// Delivery Report
	router.GET(constant.DeliveryReport, pubsub.DeliveryReportHandler)

	// start server
	log.Fatal(router.Run(common.GoDotEnvVariable(common.DotEnvVarWsServerPort)))
}

// initGin to initialise Gin network module
func initGin() {
	gin.SetMode(common.GoDotEnvVariable(common.DotEnvVarGinMode))
	router = gin.Default()
}
