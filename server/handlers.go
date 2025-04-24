package server

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"
)

// BroadcastMessage sends a formatted message to all connected clients.
func (s *Server) BroadcastMessage(msg Message) {
	formatted := FormatChatMessage(msg)
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for _, client := range s.Clients {
		client.Conn.Write([]byte(formatted))
	}
}

// ReadClientInput reads and processes client messages until disconnect.
func (s *Server) ReadClientInput(client *Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf(ColorBrightRed+"Error reading message from %s: %v"+ColorReset, client.Name, err)
			break
		}
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}

		// Handle rename command
		if strings.HasPrefix(text, "/name ") {
			newName := strings.TrimSpace(strings.TrimPrefix(text, "/name "))
			if newName == "" {
				client.Conn.Write([]byte(ColorBrightRed + "Name cannot be empty\n" + ColorReset))
				continue
			}

			// Check for name collision
			nameExists := false
			s.Mu.RLock()
			for _, c := range s.Clients {
				if c.Name == newName && c != client {
					nameExists = true
					break
				}
			}
			s.Mu.RUnlock()

			if nameExists {
				client.Conn.Write([]byte(ColorBrightRed + "Name already taken. Please try another name.\n" + ColorReset))
				continue
			}

			old := client.Name
			s.Mu.Lock()
			client.Name = newName
			s.Mu.Unlock()
			changeMsg := Message{Timestamp: time.Now().Format("2006-01-02 15:04:05"), Content: fmt.Sprintf("%s changed name to %s", old, newName)}
			s.Mu.Lock()
			s.Messages = append(s.Messages, changeMsg)
			s.Mu.Unlock()
			s.BroadcastMessage(changeMsg)
			log.Printf("Client %s changed name to %s", old, newName)
			continue
		}

		// Broadcast chat message
		msg := Message{Timestamp: time.Now().Format("2006-01-02 15:04:05"), Username: client.Name, Content: text}
		s.Mu.Lock()
		s.Messages = append(s.Messages, msg)
		s.Mu.Unlock()

		s.BroadcastMessage(msg)
		log.Printf("Message from %s: %s", client.Name, text)
	}
	// Cleanup on disconnect
	close(client.Messages)
	s.ClientChan <- client
}
