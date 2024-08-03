package utility

import "github.com/gorilla/websocket"

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	Conn     *websocket.Conn
	WsServer *WsServer
	Send     chan []byte
}

type WsServer struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		Conn:     conn,
		WsServer: wsServer,
		Send:     make(chan []byte, 256),
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

		case message := <-server.Broadcast:
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.Send <- message
	}
}
