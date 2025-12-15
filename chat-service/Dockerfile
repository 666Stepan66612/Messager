FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

#possible chat-service version
#chat-service/
#├── cmd/
#│   └── main.go
#├── internal/
#│   ├── handler/
#│   │   ├── message_handler.go
#│   │   └── chat_handler.go
#│   ├── service/
#│   │   ├── message_service.go
#│   │   └── encryption_service.go
#│   ├── repository/
#│   │   ├── message_repo.go
#│   │   └── chat_repo.go
#│   ├── models/
#│   │   ├── message.go
#│   │   └── chat.go
#│   └── config/
#│       └── config.go
#├── Dockerfile
#├── go.sum
#└── go.mod