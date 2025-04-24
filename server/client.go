package server

import (
	"net"
)

// Client represents a participant in the chat.
type Client struct {
	Conn     net.Conn
	Name     string
	Messages chan Message
}
