package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net-cat/server"
	"os"
	"strconv"
)

const (
	defaultPort = 8989
	maxClients  = 10
	usageMsg    = server.ColorBrightRed + "[USAGE]: ./TCPChat $port" + server.ColorReset
	fullMsg     = server.ColorBrightRed + "Chat is full. Please try again later.\n" + server.ColorReset
)

func main() {
	port := defaultPort
	useTUI := flag.Bool("tui", false, "Use Terminal UI")
	flag.Parse()

	// Handle port argument if provided
	args := flag.Args()
	if len(args) > 0 {
		if len(args) > 1 {
			fmt.Println(usageMsg)
			os.Exit(1)
		}
		p, err := strconv.Atoi(args[0])
		if err != nil || p < 1 || p > 65535 {
			fmt.Println(usageMsg)
			os.Exit(1)
		}
		port = p
	}

	if *useTUI {
		// Run as TUI client
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatal(err)
		}

		client, err := NewTUIClient(conn)
		if err != nil {
			log.Fatal(err)
		}

		if err := client.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		// Run as server
		server := &server.Server{
			Clients:    make(map[net.Conn]*server.Client),
			Messages:   make([]server.Message, 0),
			ClientChan: make(chan *server.Client),
		}

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		defer listener.Close()

		fmt.Printf("Listening on the port :%d\n", port)

		go server.MonitorDisconnections()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}

			server.Mu.RLock()
			nClients := len(server.Clients)
			server.Mu.RUnlock()
			if nClients >= maxClients {
				conn.Write([]byte(fullMsg))
				conn.Close()
				continue
			}

			go server.AcceptClient(conn)
		}
	}
}
