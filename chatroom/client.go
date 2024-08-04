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

type WsServer struct {
	wsServers map[string]*utility.WsServer
}

var (
	newline          = []byte{'\n'}
	ctx              = context.Background()
	chatroomWsServer = newChatroomWsServer()
	upgrader         = newUpgrader()
	redisClient      = newRedisClient()
)

// newChatroomWsServer creates a new Chatroom WsServer type
func newChatroomWsServer() *WsServer {
	return &WsServer{
		wsServers: make(map[string]*utility.WsServer),
	}
}

// newUpgrader creates a new websocket Upgrader
func newUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
}

// newRedisClient creates a new Redis Client
func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// WsHandler returns custom Chatroom gin.HandlerFunc. Create/Get WsServer w.r.t chatroom_id
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

// ServeWs handles websocket requests of a chatroom from clients requests.
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

// disconnect closes connection and clears client
func disconnect(chatroomID string, client *utility.Client) {
	unregisterFromServer(chatroomID, client)
	err := client.Conn.Close()
	if err != nil {
		return
	}
}

// unregisterFromServer to remove client from WsServer and delete WsServer from ChatroomWsServer if last client left
func unregisterFromServer(chatroomID string, client *utility.Client) {
	client.WsServer.Unregister <- client
	if client.WsServer.GetClientsCount() <= 1 {
		delete(chatroomWsServer.wsServers, chatroomID)
	}
}

// readPump to read incoming message from client
func readPump(chatroomID string, client *utility.Client) {
	defer func() {
		// disconnect client after exit from for loop
		disconnect(chatroomID, client)
	}()

	//client.conn.SetReadLimit(maxMessageSize)
	// SetReadDeadline to time.Now() + PongWait (which is < PingPeriod)
	err := client.Conn.SetReadDeadline(time.Now().Add(utility.PongWait))
	if err != nil {
		return
	}
	// set SetPongHandler to read incoming "ping" message through ticker. Used to increase SetReadDeadline
	client.Conn.SetPongHandler(func(string) error {
		// update SetReadDeadline to time.Now() + PongWait
		err := client.Conn.SetReadDeadline(time.Now().Add(utility.PongWait))
		if err != nil {
			return err
		}
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		// publish jsonMessage to redis channel:<chatroomID>
		if err := redisClient.Publish(ctx, chatroomID, jsonMessage).Err(); err != nil {
			log.Println("Publish error:", err)
			return
		}
	}
}

// writePump to send message from server to client
func writePump(chatroomID string, client *utility.Client) {
	// subscribe to redis channel:<chatroomID>
	sub := redisClient.Subscribe(ctx, chatroomID)
	// start ticker at regular interval of PingPeriod
	ticker := time.NewTicker(utility.PingPeriod)
	defer func() {
		// stop ticker
		ticker.Stop()
		// close client connection only
		err := client.Conn.Close()
		if err != nil {
			return
		}
		// stop listening to redis channel:<chatroomID>
		err = sub.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case message, ok := <-sub.Channel():
			// SetWriteDeadline to time.Now() + WriteWait
			err := client.Conn.SetWriteDeadline(time.Now().Add(utility.WriteWait))
			if err != nil {
				return
			}
			if !ok {
				// The WsServer closed the channel.
				err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			// Create NextWriter of type websocket.TextMessage
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write([]byte(message.Payload))
			if err != nil {
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// At regular interval of PingPeriod, update SetWriteDeadline to time.Now() + WriteWait
			err := client.Conn.SetWriteDeadline(time.Now().Add(utility.WriteWait))
			if err != nil {
				return
			}
			// Write message of type websocket.PingMessage
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
