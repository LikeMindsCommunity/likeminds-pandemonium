package chatroom

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/utility"
	"log"
	"net/http"
	"time"
)

type ChatroomWsServer struct {
	wsServers map[string]*utility.WsServer
}

var chatroomWsServer = newChatroomWsServer()

var (
	newline = []byte{'\n'}
	ctx     = context.Background()
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

// newChatroomWsServer creates a new WsServer type
func newChatroomWsServer() *ChatroomWsServer {
	return &ChatroomWsServer{
		wsServers: make(map[string]*utility.WsServer),
	}
}

func WsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wsServer *utility.WsServer
		chatroomID := c.Query("chatroom_id")
		if _, ok := chatroomWsServer.wsServers[chatroomID]; ok {
			wsServer = chatroomWsServer.wsServers[chatroomID]
		} else {
			wsServer = utility.NewWebsocketServer()
			go wsServer.Run()
			chatroomWsServer.wsServers[chatroomID] = wsServer
		}
		ServeWs(chatroomID, wsServer, c.Writer, c.Request)
	}
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(chatroomID string, wsServer *utility.WsServer, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := utility.NewClient(conn, wsServer)

	go writePump(chatroomID, client)
	go readPump(chatroomID, client)

	wsServer.Register <- client
}

func disconnect(chatroomID string, client *utility.Client) {
	unregisterFromServer(chatroomID, client)
	close(client.Send)
	client.Conn.Close()
}

func unregisterFromServer(chatroomID string, client *utility.Client) {
	client.WsServer.Unregister <- client
	if client.WsServer.GetClientsCount() <= 1 {
		delete(chatroomWsServer.wsServers, chatroomID)
	}
}

func readPump(chatroomID string, client *utility.Client) {
	defer func() {
		disconnect(chatroomID, client)
	}()

	//client.conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(utility.PongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(utility.PongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		if err := redisClient.Publish(ctx, chatroomID, jsonMessage).Err(); err != nil {
			log.Println("Publish error:", err)
			return
		}
	}
}

func writePump(chatroomID string, client *utility.Client) {
	sub := redisClient.Subscribe(ctx, chatroomID)
	ticker := time.NewTicker(utility.PingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
		sub.Close()
	}()
	for {
		select {
		case message, ok := <-sub.Channel():
			client.Conn.SetWriteDeadline(time.Now().Add(utility.WriteWait))
			if !ok {
				// The WsServer closed the channel.
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(message.Payload))

			// Attach queued chat messages to the current websocket message.
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(utility.WriteWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
