#!/usr/bin/env bash

COLOR_YELLOW='\033[38;5;220m'   
NEON_PINK='\033[38;5;198m'   
RESET='\033[0m'               

echo -e "${NEON_PINK}Building TCPChat...${RESET}"
go build -o TCPChat 
echo -e "${COLOR_YELLOW}Build complete. Run with './TCPChat' or './TCPChat <port>'${RESET}"