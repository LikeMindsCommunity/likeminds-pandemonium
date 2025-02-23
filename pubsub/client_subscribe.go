package pubsub

import (
	"encoding/json"
	"fmt"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/handlers"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/ws"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
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
			UUID := c.GetHeader(constant.HeadersMemberID)
			apiKey := c.GetHeader(constant.HeadersApiKey)
			deviceID := c.GetHeader(constant.HeadersDeviceID)
			sdkSource := c.GetHeader(constant.HeadersSDKSource)
			platformCode := c.GetHeader(constant.HeadersPlatformCode)
			versionCode, err := strconv.Atoi(c.GetHeader(constant.HeadersVersionCode))
			if err != nil {
				log.Printf("failed to parse version code, version code=%s, err=%s", c.GetHeader(constant.HeadersVersionCode), err)
			}
			apiVersion, err := strconv.Atoi(c.GetHeader(constant.HeadersApiVersion))
			if err != nil {
				log.Printf("failed to parse api version, api version=%s, err=%s", c.GetHeader(constant.HeadersApiVersion), err)
			}

			var chatroomID string
			if len(topicSplit) > 1 {
				chatroomID = topicSplit[1]
			}
			//Close connection if UUID sent in headers is invalid when subscribed to TopicTypeChatroom
			if UUID == "" || UUID == "null" {
				api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
				return
			}
			//Close connection if chatroom ID sent in request is invalid when subscribed to TopicTypeChatroom
			if chatroomID == "" || chatroomID == "null" {
				api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
				return
			}

			err = ServeWs(wsServerParent, topic, UUID, apiKey, deviceID, c.Writer, c.Request, redisClient, sdkSource, platformCode, versionCode, apiVersion)
			if err != nil {
				updatedErr := fmt.Sprintf(common.ErrorFailedUpgrader, err)
				api.GeneralAPIError(c, updatedErr)
				return
			}
		case common.TopicTypeCommunity:
			UUID := c.GetHeader(constant.HeadersMemberID)
			apiKey := c.GetHeader(constant.HeadersApiKey)
			deviceID := c.GetHeader(constant.HeadersDeviceID)
			sdkSource := c.GetHeader(constant.HeadersSDKSource)
			platformCode := c.GetHeader(constant.HeadersPlatformCode)
			versionCode, err := strconv.Atoi(c.GetHeader(constant.HeadersVersionCode))
			if err != nil {
				log.Printf("failed to parse version code, version code=%s, err=%s", c.GetHeader(constant.HeadersVersionCode), err)
			}
			apiVersion, err := strconv.Atoi(c.GetHeader(constant.HeadersApiVersion))
			if err != nil {
				log.Printf("failed to parse api version, api version=%s, err=%s", c.GetHeader(constant.HeadersApiVersion), err)
			}

			var communityID string
			if len(topicSplit) > 1 {
				communityID = topicSplit[1]
			}
			//Close connection if UUID sent in headers is invalid when subscribed to TopicTypeChatroom
			if UUID == "" || UUID == "null" {
				api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
				return
			}
			//Close connection if community ID sent in request is invalid when subscribed to TopicTypeChatroom
			if communityID == "" || communityID == "null" {
				api.GeneralBadRequestError(c, common.ErrorCommunityIDMissing)
				return
			}
			err = ServeWs(wsServerParent, topic, UUID, apiKey, deviceID, c.Writer, c.Request, redisClient, sdkSource, platformCode, versionCode, apiVersion)
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
		wsServerParent.WsServers[topic] = wsServer
	}
	return wsServer
}

// ServeWs handles websocket requests of a chatroom from clients requests.
func ServeWs(wsServerParent *ws.WsServerParent, topic string, UUID string, apiKey string, deviceID string, w http.ResponseWriter, r *http.Request, redisClient *redis.Client, sdkSource string, platformCode string, versionCode int, apiVersion int) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	wsServer := createOrGetWsServer(wsServerParent, topic)
	client := ws.NewClient(conn, wsServer, UUID, apiKey, deviceID, topic, sdkSource, platformCode, versionCode, apiVersion)

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
		readMessageType, readMessagePayload, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf(common.ErrorReadClientWs, err)
			return
		}
		log.Printf(common.ReceivedMessageClientWs, readMessageType)

		var readMessageJsonMap map[string]interface{}
		err = json.Unmarshal(readMessagePayload, &readMessageJsonMap)
		if err != nil {
			log.Printf(common.ErrorInvalidJSONFormat, err)
			return
		}

		topic := client.Topic
		topicSplit, err := GetTopicSplit(topic)
		if err != nil {
			log.Printf(common.ErrorTopicInvalid, err)
			return
		}

		topicMessageType := readMessageJsonMap[common.ParamTopicMessageType]

		switch topicSplit[0] {
		case common.TopicTypeChatroom:
			switch topicMessageType {
			case common.TopicMessageTypeCreateConversationRequest:
				log.Println(topicMessageType)

				participants := readMessageJsonMap[common.ParamParticipantsType]
				participantsStringList, ok := participants.([]string)
				if !ok {
					log.Print(common.ErrorInvalidTotalParticipantsFormat)
					return
				}

				totalParticipantsCount := readMessageJsonMap[common.ParamTotalParticipantsCountType]
				totalParticipantsCountInt, ok := totalParticipantsCount.(int)
				if !ok {
					log.Print(common.ErrorInvalidTotalParticipantsFormat)
					return
				}

				// Create conversation data in database and return PSResponse with topic_message_type as message.create.response
				createMessagePSResponse, createMessageResponse := handlers.CreateMessage(readMessageJsonMap, client.UUID, client.ApiKey, client.DeviceID, client.Topic, client.SDKSource, client.PlatformCode, client.VersionCode, client.ApiVersion, participantsStringList, totalParticipantsCountInt)

				// publish response to pubsub TopicNameChatroom
				createMessagePSResponseBytes, err := json.Marshal(createMessagePSResponse)
				if err != nil {
					log.Printf(common.ErrorInvalidJSONFormat, err)
					return
				}

				// Send sent delivery report to sender connection
				go updateSentDROnSubscribe(redisClient, wsServerParent, client.UUID, client.DeviceID, createMessageResponse.Data.Message.CardID,
					createMessageResponse.Data.Message.CommunityID, createMessageResponse.Data.Message.ID, float64(createMessageResponse.Data.Message.CreatedAt),
					createMessageResponse.TotalParticipantsCount)
				//todo to publish to community topicPrefix as well
				if err := PublishMessageToRedis(redisClient, client.Topic, createMessagePSResponseBytes); err != nil {
					return
				}
			}
		}
	}
}

// writePump to send message from server to client
func writePump(wsServerParent *ws.WsServerParent, client *ws.Client, redisClient *redis.Client) {
	// subscribe to pubsub TopicNameChatroom
	sub, err := SubscribeToRedisTopic(redisClient, client.Topic)
	if err != nil {
		return
	}
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
		case channelMessage, ok := <-sub.Channel():
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
			// Unmarshal channelMessageByte PSResponse
			channelMessageByte := []byte(channelMessage.Payload)
			var channelMessagePSResponse requestresponse.PSResponse
			if err := json.Unmarshal(channelMessageByte, &channelMessagePSResponse); err != nil {
				log.Printf(common.ErrorUnmarshalErrorJson, err)
				return
			}

			switch channelMessagePSResponse.TopicMessageType {
			case common.TopicMessageTypeCreateConversationResponse:
				var createMessageResponse requestresponse.CreateMessageResponse
				if err := json.Unmarshal([]byte(channelMessagePSResponse.RawData), &createMessageResponse); err != nil {
					log.Printf(common.ErrorUnmarshalErrorJson, err)
					return
				}
				// To not return to user who has sent the channelMessage and is on the same device. If user opts to not send device_id then we will send it to the same user as well
				if (*createMessageResponse.Data.User.UserUniqueID == client.UUID) &&
					(client.DeviceID != "" && client.DeviceID == channelMessagePSResponse.DeviceID) {
					continue
				}
				participants := createMessageResponse.Participants
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
				_, err = w.Write(channelMessageByte)
				if err != nil {
					log.Printf(common.ErrorUnableToWriteWs, err)
					return
				}
				if err := w.Close(); err != nil {
					log.Printf(common.ErrorWriterCloseWs, err)
					return
				}

				log.Println(common.ReceivedMessageRedisWs)
				deliveredUUID := client.UUID
				deliveredDeviceID := client.DeviceID
				communityID := createMessageResponse.Data.Message.CommunityID
				chatroomID := createMessageResponse.Data.Message.CardID
				conversationID := createMessageResponse.Data.Message.ID

				go updateDeliveredDROnSubscribe(redisClient, wsServerParent, deliveredUUID, deliveredDeviceID, communityID, chatroomID, conversationID)
				go updateReadDROnSubscribe(redisClient, wsServerParent, deliveredUUID, deliveredDeviceID, communityID, chatroomID, conversationID)
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

func updateDeliveredDROnSubscribe(redisClient *redis.Client, wsServerParent *ws.WsServerParent, deliveredUUID, deliveredDeviceID string, communityID int, chatroomID int, conversationID int64) {
	// Fetch and update the dr_conversation_<conversation_id> with delivered report
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)
	if err := UpdateDeliveredDR(redisClient, wsServerParent, deliveredUUID, deliveredDeviceID, communityID, chatroomID, conversationKey); err != nil {
		log.Println(err)
	}
}

func updateReadDROnSubscribe(redisClient *redis.Client, wsServerParent *ws.WsServerParent, readUUID, readDeviceID string, communityID int, chatroomID int, conversationID int64) {
	// Fetch and update the dr_conversation_<conversation_id>
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)
	if err := UpdateReadDR(redisClient, wsServerParent, readUUID, readDeviceID, communityID, chatroomID, conversationKey); err != nil {
		log.Println(err)
	}
}

func updateSentDROnSubscribe(redisClient *redis.Client, wsServerParent *ws.WsServerParent, senderUUID, senderDeviceID string, chatroomID, communityID int, conversationID int64, conversationCreatedAt float64, participantsCount int) {
	if err := UpdateSentDR(redisClient, wsServerParent, senderUUID, senderDeviceID, chatroomID, communityID, conversationID, conversationCreatedAt, participantsCount); err != nil {
		log.Println(err)
	}
}
