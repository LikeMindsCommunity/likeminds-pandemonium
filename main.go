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

	//wsServer := message.NewWebsocketServer()
	//go wsServer.Run()
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	message.ServeWs(wsServer, w, r)
	//})

	router.GET(utility.ChatroomListen, chatroom.WsHandler())
	log.Fatal(router.Run(":8080"))
}

func initGin() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}
