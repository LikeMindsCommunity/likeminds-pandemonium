package ws

import (
	"encoding/json"
	"fmt"
	"likeminds-pandemonium/common"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Client represents the websocket client connected with the server
type Client struct {
	// The actual websocket connection.
	Conn *websocket.Conn
	// Server to which client is connected
	WsServer *WsServer
	// UUID of user
	UUID string
	// Device ID of user
	DeviceID string
	// Topic on which client is connected
	Topic string
	// Platform Code of user
	PlatformCode string
}

// WsServer represents web socket server
type WsServer struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}

type WsServerParent struct {
	WsServers map[string]*WsServer
}

// NewClient creates new client which will be added to WsServer
func NewClient(conn *websocket.Conn, wsServer *WsServer, UUID string, deviceID string, topic string, platformCode string) *Client {
	return &Client{
		Conn:         conn,
		WsServer:     wsServer,
		UUID:         UUID,
		DeviceID:     deviceID,
		Topic:        topic,
		PlatformCode: platformCode,
	}
}

// SendPayloadToClientConnection sends the payload to the client's WebSocket connection
func (client *Client) SendPayloadToClientConnection(payload interface{}) error {
	// Marshal the payload into JSON bytes
	messagePayloadByte, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}

	// Create NextWriter of type websocket.TextMessage
	w, err := client.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf(common.ErrorWriterOpenWs, err)
	}

	// Write the payload to the WebSocket connection
	_, err = w.Write(messagePayloadByte)
	if err != nil {
		return fmt.Errorf(common.ErrorUnableToWriteWs, err)
	}

	// Close the writer
	if err := w.Close(); err != nil {
		return fmt.Errorf(common.ErrorWriterCloseWs, err)
	}

	return nil
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.Register:
			server.registerClient(client)
		case client := <-server.Unregister:
			server.unregisterClient(client)
		}
	}
}

// registerClient to register client on WsServer
func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

// unregisterClient to unregister client on WsServer
func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

// GetClientsCount to return clients count connected on WsServer
func (server *WsServer) GetClientsCount() int {
	return len(server.clients)
}

// FindClientByUUID searches through the clients map and returns the client where UUID == senderUUID
func (server *WsServer) FindClientByUUID(senderUUID string) (*Client, bool) {
	// Loop through the clients in the map
	for client := range server.clients {
		if client.UUID == senderUUID {
			return client, true // Return the matched client
		}
	}
	return nil, false // Return nil and false if no match found
}

// NewWsServerParent creates a new WsServerParent
func NewWsServerParent() *WsServerParent {
	return &WsServerParent{
		WsServers: make(map[string]*WsServer),
	}
}

// GetWsServerParentFromContext Exposed api method to get WsServerParent from context
func GetWsServerParentFromContext(c *gin.Context) *WsServerParent {
	wsServerParent, exists := c.Get(common.WsServerKey)
	if !exists {
		return nil
	}
	return wsServerParent.(*WsServerParent)
}

// GetWsServer returns the WsServer for the given topic from WsServerParent
func (parent *WsServerParent) GetWsServer(topic string) *WsServer {
	// Return the WsServer associated with the topic, or nil if not found
	return parent.WsServers[topic]
}

// GetConnectionFromWsServer function to get WebSocket connection from the server
func (parent *WsServerParent) GetConnectionFromWsServer(topic string, senderUUID string) *Client {
	wsServer := parent.GetWsServer(topic)
	if wsServer != nil {
		wsServerClient, found := wsServer.FindClientByUUID(senderUUID)
		if found {
			return wsServerClient
		}
	}
	return nil
}
