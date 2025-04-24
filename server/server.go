package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type Server struct {
	Clients    map[net.Conn]*Client
	Messages   []Message
	Mu         sync.RWMutex
	ClientChan chan *Client
}

// init initializes logging
func init() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	log.SetOutput(logFile)
}
