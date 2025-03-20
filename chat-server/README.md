- A **Chat Server** that supports real-time messaging with group chat functionality.
- Implements concepts/technologies like:
    - Real-Time Messaging : Use of WebSockets for bidirectional communication.
    - Group Chat : Allows users to join specific chat rooms (groups).
    - Message History Caching : Caching recent messages in memory for quick retrieval.
    - Pub/Sub for Scalability : Redis for publish/subscribe to handle message broadcasting across multiple server instances.
    - Database Storage : Persisting all chat messages in MongoDB using GORM as the ORM.
    - REST API : Gin to fetch chat history from the database.
    - Concurrency & State Management : goroutines and channels for managing WebSocket connections and state.

```bash
chat-server/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── message.go
│   ├── repository/
│   │   └── message_repository.go
│   ├── services/
│   │   └── chat_service.go
│   └── websocket/
│       └── websocket.go
├── pkg/
│   ├── cache/
│   │   └── cache.go
│   ├── db/
│   │   └── db.go
│   └── pubsub/
│       └── pubsub.go
├── go.mod
└── go.sum
```