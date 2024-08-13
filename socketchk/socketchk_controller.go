package socketchk

import (
	"io"
	"log"

	"github.com/NateshR/Likeminds-Real-Time/middleware"
	"github.com/gin-gonic/gin"
)

// GetSocketStatus reads and responds to socket connection request
func GetSocketStatus(c *gin.Context) {
	socketConn, ok := middleware.FromWebSocketContext(c.Request.Context())
	if !ok {
		log.Println("Failed to get websocket connection.")
		return
	}

	for {
		messageType, r, err := socketConn.NextReader()
		if err != nil {
			return
		}
		w, err := socketConn.NextWriter(messageType)
		if err != nil {
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			return
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}
