package handler

/*beta handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
	"chat-service/internal/models"
)

type MessageHandler struct{
	messageService MessageService
	authService    authServiceClient
}

func NewMessageHandler(msgSvc MessageService, authSvc AuthServiceClient) *MessageHandler{
	return &MessageHandler{
		messageService: msgSvc,
		authService:    authSvc,
	}
}

type SendMessageRequest struct{
	ChatID          string `json:"chat_id binding:"required,uuid"`
	EphemeralPubKey []byte `json:"epheral_pub_key" binding:"required"`
	Ciphertext      []byte `json:"ciphertext" binding:"required"`
	Nonce           []byte `json:"nonce" binding:"requider,len=12"`
}

type SendMessageResponse struct{
	MessageID string `json:message_id`
	Timestamp int64  `json:timestamp`
}

//
func (h *MessageHandler)SendMessage(c *gin.Context){
	//fetch and validate JWT
	token := c.GetHeader("Authorization")
	if len(token) < 8 || token[:7] != "Bearer "{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})

		return
	}

	//check token using auth-service (gRPC call)
	userID, err := h.authService.VerifyToken(token[7:])
	if err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})

		return
	}

	//parse request body
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})

		return
	}

	//validate size of ciphertext
	if len(req.Ciphertext) > 10*1024*1024{ //max 10MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "message too large"})

		return
	}

	//check that user is chat member
	isMember, err := h.messageService.IsChatMember(req.ChatID, userID)
	if err != nil || !isMember{
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this chat"})

		return
	}

	//save encoded message to db
	messageID := uuid.New().String()
	timestamp, err := h.messageService.SaveMessage(&Message{
		ID:              messageID,
        ChatID:          req.ChatID,
        SenderID:        userID,
        EphemeralPubKey: req.EphemeralPubKey,
        Ciphertext:      req.Ciphertext,
        Nonce:           req.Nonce,
        Status:          "delivered",
	})

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})

		return
	}

	go h.messageService.NotifyChatMember(req.ChatID, userID, messageID)

	c.JSON(http.StatusOK, SendMessageResponse{
		MessageID: messageID,
		Timestamp: timestamp,
	})
}

func (h *MessageHandler)GetMessages(c *gin.Context){}

func (h *MessageHandler)DeleteMessage(c *gin.Context){}

func (h *MessageHandler)MarkAsRead(c *gin.Context){}
*/