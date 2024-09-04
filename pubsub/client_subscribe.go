package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common/models"
	"likeminds-pandemonium/ws"
	"log"
	"net/http"
	"time"
)

type SubscribeWsServer struct {
	wsServers map[string]*ws.WsServer
}

var (
	subscribeWsServer = newSubscribeWsServer()
	upgrader          = newUpgrader()
)

// newSubscribeWsServer creates a new Chatroom SubscribeWsServer type
func newSubscribeWsServer() *SubscribeWsServer {
	return &SubscribeWsServer{
		wsServers: make(map[string]*ws.WsServer),
	}
}

// newUpgrader creates a new websocket Upgrader
func newUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  ReadBufferSizeDefault,
		WriteBufferSize: WriteBufferSizeDefault,
	}
}

// Subscribe to a specific topic
func Subscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		topic := c.Param(ParamTopic)
		topicSplit, err := GetTopicSplit(topic)
		if err != nil {
			api.GeneralBadRequestError(c, err.Error())
			return
		}

		switch topicSplit[0] {
		case TopicTypeChatroom:
			UUID := c.GetHeader(constant.HeadersMemberId)
			deviceID := c.GetHeader(constant.HeadersDeviceId)
			var chatroomID string
			if len(topicSplit) > 1 {
				chatroomID = topicSplit[1]
			}
			//Close connection if UUID sent in headers is invalid when subscribed to TopicTypeChatroom
			if UUID == "" || UUID == "null" {
				api.GeneralBadRequestError(c, ErrorUserUUIDMissing)
				return
			}
			//Close connection if chatroom ID sent in request is invalid when subscribed to TopicTypeChatroom
			if chatroomID == "" || chatroomID == "null" {
				api.GeneralBadRequestError(c, ErrorChatroomIDMissing)
				return
			}

			err = ServeWs(topic, UUID, deviceID, createOrGetWsServer(topic), c.Writer, c.Request, GetRedisClientFromContext(c))
			if err != nil {
				updatedErr := fmt.Sprintf(ErrorFailedUpgrader, err)
				api.GeneralAPIError(c, updatedErr)
				return
			}
		case TopicTypeCommunity:
			UUID := c.GetHeader(constant.HeadersMemberId)
			deviceID := c.GetHeader(constant.HeadersDeviceId)
			var communityID string
			if len(topicSplit) > 1 {
				communityID = topicSplit[1]
			}
			//Close connection if UUID sent in headers is invalid when subscribed to TopicTypeChatroom
			if UUID == "" || UUID == "null" {
				api.GeneralBadRequestError(c, ErrorUserUUIDMissing)
				return
			}
			//Close connection if community ID sent in request is invalid when subscribed to TopicTypeChatroom
			if communityID == "" || communityID == "null" {
				api.GeneralBadRequestError(c, ErrorCommunityIDMissing)
				return
			}
			wsServer := createOrGetWsServer(topic)
			redisClient := GetRedisClientFromContext(c)
			err = ServeWs(topic, UUID, deviceID, wsServer, c.Writer, c.Request, redisClient)
			if err != nil {
				updatedErr := fmt.Sprintf(ErrorFailedUpgrader, err)
				api.GeneralAPIError(c, updatedErr)
				return
			}
		}
	}
}

// createOrGetWsServer to create new SubscribeWsServer get against `topic`
func createOrGetWsServer(topic string) *ws.WsServer {
	var wsServer *ws.WsServer
	if _, ok := subscribeWsServer.wsServers[topic]; ok {
		wsServer = subscribeWsServer.wsServers[topic]
	} else {
		wsServer = ws.NewWebsocketServer()
		go wsServer.Run()
	}
	return wsServer
}

// ServeWs handles websocket requests of a chatroom from clients requests.
func ServeWs(topic string, UUID string, deviceID string, wsServer *ws.WsServer, w http.ResponseWriter, r *http.Request, redisClient *redis.Client) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := ws.NewClient(conn, wsServer, UUID, deviceID, topic)

	go writePump(client, redisClient)
	go readPump(client, redisClient)

	wsServer.Register <- client
	log.Println(WsConnectionEstablished)
	return nil
}

// disconnect closes connection and clears client
func disconnect(client *ws.Client) {
	unregisterFromServer(client)
	err := client.Conn.Close()
	if err != nil {
		log.Println(ErrorUnableToCloseWs, err)
		return
	}
}

// unregisterFromServer to remove client from SubscribeWsServer and delete SubscribeWsServer from ChatroomWsServer if last client left
func unregisterFromServer(client *ws.Client) {
	client.WsServer.Unregister <- client
	if client.WsServer.GetClientsCount() <= 1 {
		delete(subscribeWsServer.wsServers, client.Topic)
	}
}

// readPump to read incoming message from client
func readPump(client *ws.Client, redisClient *redis.Client) {
	defer func() {
		// disconnect client after exit from for loop
		disconnect(client)
	}()

	//client.conn.SetReadLimit(maxMessageSize)
	// SetReadDeadline to time.Now() + PongWait (which is < PingPeriod)
	err := client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	if err != nil {
		log.Println(ErrorReadDeadlineWs, err)
		return
	}
	// set SetPongHandler to read incoming "ping" message through ticker. Used to increase SetReadDeadline
	client.Conn.SetPingHandler(func(string) error {
		log.Println(PingWs)
		// update SetReadDeadline to time.Now() + PongWait
		err := client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
		if err != nil {
			log.Println(ErrorReadDeadlineWs, err)
			return err
		}
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		messageType, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(fmt.Sprintf(ErrorReadClientWs, err))
			return
		}
		log.Println(fmt.Sprintf(ReceivedMessageClientWs, messageType))
		// publish jsonMessage to pubsub TopicNameChatroom
		if err := redisClient.Publish(ctx, client.Topic, jsonMessage).Err(); err != nil {
			log.Println(ErrorPublishRedis, err)
			return
		}
	}
}

// writePump to send message from server to client
func writePump(client *ws.Client, redisClient *redis.Client) {
	// subscribe to pubsub TopicNameChatroom
	sub := redisClient.Subscribe(ctx, client.Topic)
	//// start ticker at regular interval of PingPeriod
	//ticker := time.NewTicker(ws.PingPeriod)
	defer func() {
		//// stop ticker
		//ticker.Stop()
		// close client connection only
		err := client.Conn.Close()
		if err != nil {
			log.Println(ErrorUnableToCloseWs, err)
			return
		}
		// stop listening to pubsub TopicNameChatroom
		err = sub.Close()
		if err != nil {
			log.Println(ErrorUnableToCloseRedis, err)
			return
		}
	}()
	for {
		select {
		case message, ok := <-sub.Channel():
			// SetWriteDeadline to time.Now() + WriteWait
			err := client.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
			if err != nil {
				log.Println(ErrorWriteDeadlineWs, err)
				return
			}
			if !ok {
				// The SubscribeWsServer closed the channel.
				err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println(ErrorUnableToWriteWs, err)
					return
				}
				return
			}
			// Unmarshal messagePayloadByte Response
			messagePayloadByte := []byte(message.Payload)
			var response Response
			if err := json.Unmarshal(messagePayloadByte, &response); err != nil {
				log.Println(ErrorUnmarshalErrorJson, err)
				return
			}
			switch response.TopicMessageType {
			case TopicMessageTypeConversation:
				var conversationResponse models.ConversationResponse
				if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
					log.Println(ErrorUnmarshalErrorJson, err)
					return
				}
				// To not return to user who has sent the message and is on the same device. If user opts to not send device_id then we will send it to the same user as well
				if (conversationResponse.Conversation.Member.UUID == client.UUID) &&
					(client.DeviceID != "" && client.DeviceID == response.DeviceID) {
					continue
				}
				// Create NextWriter of type websocket.TextMessage
				w, err := client.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Println(ErrorWriterOpenWs, err)
					return
				}
				_, err = w.Write(messagePayloadByte)
				if err != nil {
					log.Println(ErrorUnableToWriteWs, err)
					return
				}
				if err := w.Close(); err != nil {
					log.Println(ErrorWriterCloseWs, err)
					return
				}
				log.Println(ReceivedMessageRedisWs)
			}
			/*case <-ticker.C:
			// At regular interval of PingPeriod, update SetWriteDeadline to time.Now() + WriteWait
			t := time.Now().Add(ws.WriteWait)
			err := client.Conn.SetWriteDeadline(t)
			log.Println("sent ping")
			if err != nil {
				log.Println(ErrorWriteDeadlineWs, err)
				return
			}
			// Write message of type websocket.PingMessage
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(ErrorUnableToWriteWs, err)
				return
			}*/
		}
	}
}
