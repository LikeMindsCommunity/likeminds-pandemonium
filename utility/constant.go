package utility

import "time"

const (
	// WriteWait Max wait time when writing message to peer
	WriteWait = 10 * time.Second

	// PongWait Max time till next pong from peer
	PongWait = 60 * time.Second

	// PingPeriod Send ping interval, must be less then pong wait time
	PingPeriod = (PongWait * 9) / 10
)
const (
	ChatroomListen = "/chatroom/listen"
)
