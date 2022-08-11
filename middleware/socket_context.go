package middleware

import (
	"context"

	"github.com/gorilla/websocket"
)

type key int

const socketConnKey key = 0

// NewWebSocketContext return a new context that carries websocket connection
func NewWebSocketContext(ctx context.Context, conn *websocket.Conn) context.Context {
	return context.WithValue(ctx, socketConnKey, conn)
}

// FromWebSocketContext extracts websocket connection from a context
func FromWebSocketContext(ctx context.Context) (*websocket.Conn, bool) {
	conn, ok := ctx.Value(socketConnKey).(*websocket.Conn)
	return conn, ok
}
