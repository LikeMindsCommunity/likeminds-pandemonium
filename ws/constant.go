package ws

import "time"

const (
	// WriteWait Max wait time when writing message to peer
	WriteWait = 10 * time.Second

	// PongWait Max time till next pong from peer
	PongWait = 60 * time.Second

	// PingPeriod should be less than PongWait
	PingPeriod = ((PongWait * time.Second) * 9) / 10
)
