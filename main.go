package main

import (
	"github.com/gin-gonic/gin"
	"likeminds-pandemonium/chatroom"
	"likeminds-pandemonium/utility"
	"log"
)

var (
	router *gin.Engine
)

func main() {
	initGin()

	// ChatroomListen GET request
	router.GET(utility.ChatroomListen, chatroom.WsHandler())

	// start server
	log.Fatal(router.Run(":8080"))
}

// initGin to initialise Gin network module
func initGin() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}
