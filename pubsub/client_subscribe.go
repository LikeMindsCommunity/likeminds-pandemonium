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

	go writePump(wsServerParent, client, redisClient, topic)
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
			log.Printf(common.ErrorReadClientWs, err)
			return
		}
		log.Printf(common.ReceivedMessageClientWs, messageType)

		var jsonMessageMap map[string]interface{}
		err = json.Unmarshal(jsonMessage, &jsonMessageMap)
		if err != nil {
			log.Printf(common.ErrorInvalidJSONFormat, err)
			return
		}
		topicMessageType := jsonMessageMap[common.ParamTopicMessageType]
		switch topicMessageType {
		case common.TopicMessageTypeCreateConversationRequest:
			log.Println(topicMessageType)

			participants := jsonMessageMap[common.ParamParticipantsType]
			participantsStringList, ok := participants.([]string)
			if !ok {
				log.Print(common.ErrorInvalidTotalParticipantsFormat)
				return
			}

			totalParticipantsCount := jsonMessageMap[common.ParamTotalParticipantsCountType]
			totalParticipantsCountInt, ok := totalParticipantsCount.(int)
			if !ok {
				log.Print(common.ErrorInvalidTotalParticipantsFormat)
				return
			}

			// Create conversation data in database
			createMessagePSResponse, createMessageResponse := handlers.CreateMessage(jsonMessageMap, client.UUID, client.ApiKey, client.DeviceID, client.Topic, client.SDKSource, client.PlatformCode, client.VersionCode, client.ApiVersion, participantsStringList, totalParticipantsCountInt)

			// publish response to pubsub TopicNameChatroom
			createConversationPSResponseBytes, err := json.Marshal(createMessagePSResponse)
			if err != nil {
				log.Printf(common.ErrorInvalidJSONFormat, err)
				return
			}

			// Send sent delivery report to sender connection
			go updateSentDR(redisClient, wsServerParent, client.UUID, client.DeviceID, createMessageResponse.Data.Message.CardID,
				createMessageResponse.Data.Message.CommunityID, createMessageResponse.Data.Message.ID, float64(createMessageResponse.Data.Message.CreatedAt),
				createMessageResponse.TotalParticipantsCount)
			//todo to publish to community topic as well
			if err := PublishMessageToRedis(redisClient, client.Topic, createConversationPSResponseBytes); err != nil {
				return
			}
		default:
			if err := PublishMessageToRedis(redisClient, client.Topic, jsonMessage); err != nil {
				return
			}
		}
	}
}

// writePump to send message from server to client
func writePump(wsServerParent *ws.WsServerParent, client *ws.Client, redisClient *redis.Client, topic string) {
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
			// Unmarshal messagePayloadByte PSResponse
			messagePayloadByte := []byte(message.Payload)
			var psResponse requestresponse.PSResponse
			if err := json.Unmarshal(messagePayloadByte, &psResponse); err != nil {
				log.Printf(common.ErrorUnmarshalErrorJson, err)
				return
			}

			switch psResponse.TopicMessageType {
			case common.TopicMessageTypeCreateConversationResponse:
				var conversationResponse requestresponse.CreateMessageResponse
				if err := json.Unmarshal([]byte(psResponse.RawData), &conversationResponse); err != nil {
					log.Printf(common.ErrorUnmarshalErrorJson, err)
					return
				}
				// To not return to user who has sent the message and is on the same device. If user opts to not send device_id then we will send it to the same user as well
				if (*conversationResponse.Data.User.UserUniqueID == client.UUID) &&
					(client.DeviceID != "" && client.DeviceID == psResponse.DeviceID) {
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
				go updateDeliveredDROnSubscribe(redisClient, wsServerParent, &conversationResponse, client.DeviceID, client.UUID)
				go updateReadDROnSubscribe(redisClient, wsServerParent, &conversationResponse, client.DeviceID, client.UUID)
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

func updateDeliveredDROnSubscribe(redisClient *redis.Client, wsServerParent *ws.WsServerParent, conversationResponse *requestresponse.CreateMessageResponse, deliveredDeviceID, deliveredUUID string) {
	if conversationResponse == nil {
		return
	}

	conversationID := conversationResponse.Data.Message.ID
	chatroomID := conversationResponse.Data.Message.CardID
	communityID := conversationResponse.Data.Message.CommunityID

	if err := UpdateDeliveredDRWithConversationID(redisClient, wsServerParent, communityID, chatroomID, conversationID, deliveredUUID, deliveredDeviceID); err != nil {
		log.Println(err)
	}
}

func updateReadDROnSubscribe(redisClient *redis.Client, wsServerParent *ws.WsServerParent, conversationResponse *requestresponse.CreateMessageResponse, deliveredDeviceID, deliveredUUID string) {
	if conversationResponse == nil {
		return
	}
	conversationID := conversationResponse.Data.Message.ID
	chatroomID := conversationResponse.Data.Message.CardID
	communityID := conversationResponse.Data.Message.CommunityID

	if err := UpdateReadDRWithConversationID(redisClient, wsServerParent, chatroomID, conversationID, deliveredDeviceID, deliveredUUID, communityID); err != nil {
		log.Println(err)
	}
}

func updateSentDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, senderUUID, senderDeviceID string, chatroomID, communityID int, conversationID int64, conversationCreatedAt float64, participantsCount int) {
	if err := UpdateSentDR(redisClient, wsServerParent, senderUUID, senderDeviceID, chatroomID, communityID, conversationID, conversationCreatedAt, participantsCount); err != nil {
		log.Println(err)
	}
}
