package server

import (
	"fmt"
	"strings"
)

type Message struct {
	Timestamp string
	Username  string
	Content   string
}

func FormatChatMessage(msg Message) string {
	if msg.Username == "" {
		if strings.Contains(msg.Content, "has left our chat...") {
			return ColorCyan + msg.Content + ColorReset + "\n"
		}
		return ColorYellow + msg.Content + ColorReset + "\n"
	}
	return fmt.Sprintf(ColorGreen+"[%s][%s]"+ColorReset+": %s\n", msg.Timestamp, msg.Username, msg.Content)
}
