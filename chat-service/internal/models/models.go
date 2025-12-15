package models

import "time"

type Message struct{
	ID              string    `json:"id" db:"id"`
    ChatID          string    `json:"chat_id" db:"chat_id"`
    SenderID        string    `json:"sender_id" db:"sender_id"`
    EphemeralPubKey []byte    `json:"ephemeral_pub_key" db:"ephemeral_pub_key"`
    Ciphertext      []byte    `json:"ciphertext" db:"ciphertext"`
    Nonce           []byte    `json:"nonce" db:"nonce"`
    Status          string    `json:"status" db:"status"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Chat struct{
	ID        string    `json:"id" db:"id"`
    Type      string    `json:"type" db:"type"`
    Name      string    `json:"name" db:"name"`
    CreatorID string    `json:"creator_id" db:"creator_id"`
    Members   []string  `json:"members" db:"members"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}