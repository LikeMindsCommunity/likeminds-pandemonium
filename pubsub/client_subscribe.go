package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/common/models"
	"likeminds-pandemonium/ws"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = newUpgrader()
)

// newUpgrader creates a new websocket Upgrader
func newUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  common.ReadBufferSizeDefault,
		WriteBufferSize: common.WriteBufferSizeDefault,
	}
}

// Subscribe to a specific topic
func Subscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		redisClient := GetRedisClientFromContext(c)
		wsServerParent := ws.GetWsServerParentFromContext(c)
		topic := c.Param(common.ParamTopic)
		topicSplit, err := GetTopicSplit(topic)
		if err != nil {
			api.GeneralBadRequestError(c, err.Error())
			return
		}

		switch topicSplit[0] {
		case common.TopicTypeChatroom:
			UUID := c.GetHeader(constant.HeadersMemberId)
			deviceID := c.GetHeader(constant.HeadersDeviceId)
			var chatroomID string
			if len(topicSplit) > 1 {
				chatroomID = topicSplit[1]
			}
			//Close connection if UUID sent in headers is invalid when subscribed to TopicTypeChatroom
			if UUID == "" || UUID == "null" {
				api.GeneralBadRequestError(c, common.ErrorUserUUIDMissing)
				return
			}
			//Close connection if chatroom ID sent in request is invalid when subscribed to TopicTypeChatroom
			if chatroomID == "" || chatroomID == "null" {
				api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
				return
			}

			err = ServeWs(wsServerParent, topic, UUID, deviceID, c.Writer, c.Request, redisClient)
			if err != nil {
				updatedErr := fmt.Sprintf(common.ErrorFailedUpgrader, err)
				api.GeneralAPIError(c, updatedErr)
				return
			}
		case common.TopicTypeCommunity:
			UUID := c.GetHeader(constant.HeadersMemberId)
			deviceID := c.GetHeader(constant.HeadersDeviceId)
			var communityID string
			if len(topicSplit) > 1 {
				communityID = topicSplit[1]
			}
			//Close connection if UUID sent in headers is invalid when subscribed to TopicTypeChatroom
			if UUID == "" || UUID == "null" {
				api.GeneralBadRequestError(c, common.ErrorUserUUIDMissing)
				return
			}
			//Close connection if community ID sent in request is invalid when subscribed to TopicTypeChatroom
			if communityID == "" || communityID == "null" {
				api.GeneralBadRequestError(c, common.ErrorCommunityIDMissing)
				return
			}
			err = ServeWs(wsServerParent, topic, UUID, deviceID, c.Writer, c.Request, redisClient)
			if err != nil {
				updatedErr := fmt.Sprintf(common.ErrorFailedUpgrader, err)
				api.GeneralAPIError(c, updatedErr)
				return
			}
		}
	}
}

// createOrGetWsServer to create new WsServerParent get against `topic`
func createOrGetWsServer(wsServerParent *ws.WsServerParent, topic string) *ws.WsServer {
	wsServer := wsServerParent.GetWsServer(topic)
	if wsServer == nil {
		wsServer = ws.NewWebsocketServer()
		go wsServer.Run()
	}
	return wsServer
}

// ServeWs handles websocket requests of a chatroom from clients requests.
func ServeWs(wsServerParent *ws.WsServerParent, topic string, UUID string, deviceID string, w http.ResponseWriter, r *http.Request, redisClient *redis.Client) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	wsServer := createOrGetWsServer(wsServerParent, topic)
	client := ws.NewClient(conn, wsServer, UUID, deviceID, topic)

	go writePump(wsServerParent, client, redisClient)
	go readPump(wsServerParent, client, redisClient)

	wsServer.Register <- client
	log.Println(common.WsConnectionEstablished)
	return nil
}

// disconnect closes connection and clears client
func disconnect(wsServerParent *ws.WsServerParent, client *ws.Client) {
	log.Println(common.ConnectionClosed)
	unregisterFromServer(wsServerParent, client)
	err := client.Conn.Close()
	if err != nil {
		log.Println(common.ErrorUnableToCloseWs, err)
		return
	}
}

// unregisterFromServer to remove client from WsServerParent and delete WsServerParent from ChatroomWsServer if last client left
func unregisterFromServer(wsServerParent *ws.WsServerParent, client *ws.Client) {
	client.WsServer.Unregister <- client
	if client.WsServer.GetClientsCount() <= 1 {
		delete(wsServerParent.WsServers, client.Topic)
	}
}

// readPump to read incoming message from client
func readPump(wsServerParent *ws.WsServerParent, client *ws.Client, redisClient *redis.Client) {
	defer func() {
		// disconnect client after exit from for loop
		disconnect(wsServerParent, client)
	}()

	updateReadDeadline(client.Conn)

	// set SetPingHandler to read incoming "ping" message. Used to increase SetReadDeadline
	client.Conn.SetPingHandler(func(string) error {
		log.Println(common.PingReceivedClient)
		sendPongMessage(client.Conn)
		updateReadDeadline(client.Conn)
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		messageType, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(fmt.Sprintf(common.ErrorReadClientWs, err))
			return
		}
		log.Println(fmt.Sprintf(common.ReceivedMessageClientWs, messageType))
		// publish jsonMessage to pubsub TopicNameChatroom
		if err := redisClient.Publish(ctx, client.Topic, jsonMessage).Err(); err != nil {
			log.Println(common.ErrorPublishRedis, err)
			return
		}
	}
}

// writePump to send message from server to client
func writePump(wsServerParent *ws.WsServerParent, client *ws.Client, redisClient *redis.Client) {
	// subscribe to pubsub TopicNameChatroom
	sub := redisClient.Subscribe(ctx, client.Topic)
	defer func() {
		disconnect(wsServerParent, client)
		// stop listening to pubsub TopicNameChatroom
		err := sub.Close()
		if err != nil {
			log.Println(common.ErrorUnableToCloseRedis, err)
			return
		}
	}()
	for {
		select {
		case message, ok := <-sub.Channel():
			updateWriteDeadline(client.Conn)
			if !ok {
				// The WsServerParent closed the channel.
				err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf(common.ErrorUnableToWriteWs, err)
					return
				}
				return
			}
			// Unmarshal messagePayloadByte Response
			messagePayloadByte := []byte(message.Payload)
			var response Response
			if err := json.Unmarshal(messagePayloadByte, &response); err != nil {
				log.Printf(common.ErrorUnmarshalErrorJson, err)
				return
			}
			switch response.TopicMessageType {
			case common.TopicMessageTypeConversation:
				var conversationResponse models.ConversationResponse
				if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
					log.Printf(common.ErrorUnmarshalErrorJson, err)
					return
				}
				// To not return to user who has sent the message and is on the same device. If user opts to not send device_id then we will send it to the same user as well
				if (conversationResponse.Conversation.Member.UUID == client.UUID) &&
					(client.DeviceID != "" && client.DeviceID == response.DeviceID) {
					continue
				}
				participants := conversationResponse.Participants
				if participants != nil && len(participants) > 0 {
					if !api.Contains(participants, client.UUID) {
						continue
					}
				}
				// Create NextWriter of type websocket.TextMessage
				w, err := client.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Printf(common.ErrorWriterOpenWs, err)
					return
				}
				_, err = w.Write(messagePayloadByte)
				if err != nil {
					log.Printf(common.ErrorUnableToWriteWs, err)
					return
				}
				if err := w.Close(); err != nil {
					log.Printf(common.ErrorWriterCloseWs, err)
					return
				}
				log.Println(common.ReceivedMessageRedisWs)
			}
		}
	}
}

func sendPongMessage(conn *websocket.Conn) {
	updateWriteDeadline(conn)
	if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
		log.Printf(fmt.Sprintf(common.ErrorPongSentClient, err))
		return
	}
	log.Println(common.PongSendClient)

}

func updateReadDeadline(conn *websocket.Conn) {
	//client.conn.SetReadLimit(maxMessageSize)
	// SetReadDeadline to time.Now() + PongWait (which is < PingPeriod)
	err := conn.SetReadDeadline(time.Now().Add(common.PongWait))
	if err != nil {
		log.Println(common.ErrorReadDeadlineWs, err)
		return
	}
}

func updateWriteDeadline(conn *websocket.Conn) {
	// SetWriteDeadline to time.Now() + WriteWait
	err := conn.SetWriteDeadline(time.Now().Add(common.WriteWait))
	if err != nil {
		log.Println(common.ErrorWriteDeadlineWs, err)
		return
	}
}
