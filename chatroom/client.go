package chatroom

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/pubsub"
	"likeminds-pandemonium/ws"
	"log"
	"net/http"
	"time"
)

type WsServer struct {
	wsServers map[string]*ws.WsServer
}

var (
	ctx              = context.Background()
	chatroomWsServer = newChatroomWsServer()
	upgrader         = newUpgrader()
)

// newChatroomWsServer creates a new Chatroom WsServer type
func newChatroomWsServer() *WsServer {
	return &WsServer{
		wsServers: make(map[string]*ws.WsServer),
	}
}

// newUpgrader creates a new websocket Upgrader
func newUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
}

// WsHandler returns custom Chatroom gin.HandlerFunc. Create/Get WsServer w.r.t chatroom_id
func WsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		chatroomID := c.Query("chatroom_id")
		UUID := c.GetHeader(constant.HeadersMemberId)
		deviceID := c.GetHeader(constant.HeadersDeviceId)
		if chatroomID == "" || chatroomID == "null" || UUID == "" || UUID == "null" || deviceID == "" || deviceID == "null" {
			return
		}
		topic := fmt.Sprintf(pubsub.TopicChatroom, chatroomID)

		var wsServer *ws.WsServer
		if _, ok := chatroomWsServer.wsServers[topic]; ok {
			wsServer = chatroomWsServer.wsServers[topic]
		} else {
			wsServer = ws.NewWebsocketServer()
			go wsServer.Run()
			chatroomWsServer.wsServers[topic] = wsServer
		}
		redisClient := pubsub.GetRedisClientFromContext(c)
		ServeWs(topic, UUID, deviceID, wsServer, c.Writer, c.Request, redisClient)
	}
}

// ServeWs handles websocket requests of a chatroom from clients requests.
func ServeWs(topic string, UUID string, deviceID string, wsServer *ws.WsServer, w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := ws.NewClient(conn, wsServer, UUID, deviceID)

	go writePump(topic, client, redisClient)
	go readPump(topic, client, redisClient)

	wsServer.Register <- client
}

// disconnect closes connection and clears client
func disconnect(topic string, client *ws.Client) {
	unregisterFromServer(topic, client)
	err := client.Conn.Close()
	if err != nil {
		return
	}
}

// unregisterFromServer to remove client from WsServer and delete WsServer from ChatroomWsServer if last client left
func unregisterFromServer(topic string, client *ws.Client) {
	client.WsServer.Unregister <- client
	if client.WsServer.GetClientsCount() <= 1 {
		delete(chatroomWsServer.wsServers, topic)
	}
}

// readPump to read incoming message from client
func readPump(topic string, client *ws.Client, redisClient *redis.Client) {
	defer func() {
		// disconnect client after exit from for loop
		disconnect(topic, client)
	}()

	//client.conn.SetReadLimit(maxMessageSize)
	// SetReadDeadline to time.Now() + PongWait (which is < PingPeriod)
	err := client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	if err != nil {
		return
	}
	// set SetPongHandler to read incoming "ping" message through ticker. Used to increase SetReadDeadline
	client.Conn.SetPongHandler(func(string) error {
		// update SetReadDeadline to time.Now() + PongWait
		err := client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
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
		// publish jsonMessage to pubsub TopicChatroom
		if err := redisClient.Publish(ctx, topic, jsonMessage).Err(); err != nil {
			log.Println("Publish error:", err)
			return
		}
	}
}

// writePump to send message from server to client
func writePump(topic string, client *ws.Client, redisClient *redis.Client) {
	// subscribe to pubsub TopicChatroom
	sub := redisClient.Subscribe(ctx, topic)
	// start ticker at regular interval of PingPeriod
	ticker := time.NewTicker(ws.PingPeriod)
	defer func() {
		// stop ticker
		ticker.Stop()
		// close client connection only
		err := client.Conn.Close()
		if err != nil {
			return
		}
		// stop listening to pubsub TopicChatroom
		err = sub.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case message, ok := <-sub.Channel():
			// SetWriteDeadline to time.Now() + WriteWait
			err := client.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
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
			err := client.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
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
