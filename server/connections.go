package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// AcceptClient handles the lifecycle of a new client connection.
func (s *Server) AcceptClient(conn net.Conn) {
	defer conn.Close()

	// Send welcome and prompt
	conn.Write([]byte(WelcomeLogo))
	conn.Write([]byte(ColorPink + "[ENTER YOUR NAME]: " + ColorReset))

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf(ColorBrightRed+"Error reading name: %v"+ColorReset, err)
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		conn.Write([]byte(ColorBrightRed + "Name cannot be empty\n" + ColorReset))
		return
	}

	// Keep trying until we get a unique name
	for {
		// Check for duplicate names
		nameExists := false
		s.Mu.RLock()
		for _, c := range s.Clients {
			if c.Name == name {
				nameExists = true
				break
			}
		}
		s.Mu.RUnlock()

		if !nameExists {
			break // Name is unique, proceed
		}

		// If name exists, ask for a new one
		conn.Write([]byte(ColorBrightRed + "Name already taken. Please try another name.\n" + ColorReset))
		conn.Write([]byte(ColorPink + "[ENTER YOUR NAME]: " + ColorReset))
		newName, err := reader.ReadString('\n')
		if err != nil {
			log.Printf(ColorBrightRed+"Error reading new name: %v"+ColorReset, err)
			return
		}
		newName = strings.TrimSpace(newName)
		if newName == "" {
			conn.Write([]byte(ColorBrightRed + "Name cannot be empty\n" + ColorReset))
			return
		}
		name = newName
	}

	client := &Client{Conn: conn, Name: name, Messages: make(chan Message, 10)}

	// Replay chat history to the new client
	s.Mu.RLock()
	for _, msg := range s.Messages {
		conn.Write([]byte(FormatChatMessage(msg)))
	}
	s.Mu.RUnlock()

	// Register the new client
	s.Mu.Lock()
	s.Clients[conn] = client
	s.Mu.Unlock()

	// Record and broadcast the join event
	joinMsg := Message{Timestamp: time.Now().Format("2006-01-02 15:04:05"), Content: fmt.Sprintf("%s has joined our chat...", client.Name)}
	s.Mu.Lock()
	s.Messages = append(s.Messages, joinMsg)
	s.Mu.Unlock()
	s.BroadcastMessage(joinMsg)
	log.Printf("Client connected: %s from %s", client.Name, conn.RemoteAddr())

	// Begin listening for client input
	go s.ReadClientInput(client)

	// Block until disconnect signal
	<-client.Messages
}

// MonitorDisconnections listens for disconnect signals and cleans up.
func (s *Server) MonitorDisconnections() {
	for client := range s.ClientChan {
		s.Mu.Lock()
		delete(s.Clients, client.Conn)
		s.Mu.Unlock()

		leaveMsg := Message{Timestamp: time.Now().Format("2006-01-02 15:04:05"), Content: fmt.Sprintf("%s has left our chat...", client.Name)}
		s.Mu.Lock()
		s.Messages = append(s.Messages, leaveMsg)
		s.Mu.Unlock()
		s.BroadcastMessage(leaveMsg)
		log.Printf("Client disconnected: %s", client.Name)
	}
}
