package ws

import "github.com/gorilla/websocket"

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
}

// WsServer represents web socket server
type WsServer struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

// NewClient creates new client which will be added to WsServer
func NewClient(conn *websocket.Conn, wsServer *WsServer, UUID string, deviceID string) *Client {
	return &Client{
		Conn:     conn,
		WsServer: wsServer,
		UUID:     UUID,
		DeviceID: deviceID,
	}
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
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
