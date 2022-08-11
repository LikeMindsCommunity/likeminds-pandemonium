package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// HttpConnectionUpgrader upgrades http connection to websocket connection
func HttpConnectionUpgrader() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(fmt.Sprintln("Error upgrading to websocket connection, reason=%s", err))
			return
		}

		ctx := NewWebSocketContext(c.Request.Context(), conn)
		c.Request = c.Request.Clone(ctx)

		c.Next()
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
