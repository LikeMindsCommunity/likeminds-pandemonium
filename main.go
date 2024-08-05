package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/chatroom"
	"likeminds-pandemonium/redisPandemonium"
	"log"
)

var (
	router      *gin.Engine
	redisClient *redis.Client
)

func main() {
	initGin()
	initRedisClient()

	router.Use(redisPandemonium.ApiMiddleware(redisClient))
	// ChatroomListen GET request
	router.GET(constant.ChatroomListen, chatroom.WsHandler())
	// Redis publish / subscribe APIs
	router.POST(constant.RedisPublish, redisPandemonium.Publish)

	// start server
	log.Fatal(router.Run(":8080"))
}

// initGin to initialise Gin network module
func initGin() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

// initRedisClient creates a new Redis Client
func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
