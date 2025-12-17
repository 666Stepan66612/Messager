package models

import "time"

type User struct{
	ID int64 `json:"id"  db:"id"`
	Username string `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}