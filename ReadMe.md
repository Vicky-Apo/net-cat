# 💬 TCP-Chat

TCP-Chat is a retro-inspired, colorful group chat server written in Go.  
Connect with friends in real time, enjoy a modern terminal experience, and relive the classic chatroom vibe!

---

## 🌟 Live Features
- 👥 **Multi-Client Support:** Up to 10 concurrent connections. Unique username required for each participant
- 🎨 **Colorful Interface:** Retro-dark mode with distinct colors for different message types and for server commands
- 📝 **Chat History:** New users see previous chat history on join
- 🔄 **Dynamic Usernames:** Change names with simple commands
- 📊 **Activity Logging:** All events saved to `server.log`


## 🗂️ Project Structure

```
.
├── main.go         # Server initialization
│   └── main()      # Entry point
├── server.go       # Core functionality
│   ├── AcceptClient()     # New connection handling
│   ├── BroadcastMessage() # Message distribution
│   └── MonitorDisconnections() # Client cleanup
├── test.go         # Unit tests
└── build.sh        # Build script
```

---

## ⚙️ How to Run

1. **Build the Server**
   ```bash
   ./build.sh
   ```

2. **Start the Server**
   ```bash
   ./TCPChat        # Default port 8989
   # or
   ./TCPChat 2525   # Custom port
   ```

3. **Connect as Client**
   ```bash
   nc localhost 8989
   ```

---

## 💡 Chat Commands & Features

### Change Username
```bash
/name NewUsername
```

### Message Format
```go
// Regular chat message
[2025-04-24 20:13:19][Username]: Message

// System message (join/leave)
Username has joined our chat...
Username has left our chat...
```

---

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

Example test case:
```go
func TestClientConnection(t *testing.T) {
    // Test client connection handling
    server := NewServer()
    // ... test implementation
}
```

---

## 👩‍💻 About the Team

We are passionate developers at Zone01, creating engaging chat experiences:

- **Vicky Apostolou, Yana Kopylova, Kostas Apostolou**

👩‍💻 Team
Created with passion by Zone01 students:

**Vicky Apostolou, Iana kopylova, Kostas Apostolou**

---

## 📁 Resources
- [Project Instructions](https://github.com/01-edu/public/tree/master/subjects/net-cat)
- [Go net Package](https://pkg.go.dev/net)
- [Go sync Package](https://pkg.go.dev/sync)

---


## © License

Made with 💚 at [Zone01](https://01.al)