package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"messenger/auth-service/internal/handlers"
    "messenger/auth-service/internal/repositories"
    "messenger/auth-service/internal/services"
)

func main(){
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal("Failed to connect to database:", err)
	}

	rdb := redis.NewClient(&redis.Option{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	db.AutoMigrate(&models.User{}, &models.Session{})

	userRepo := repositories.NewUserRepository(db)
	session := repositories.NewSessionRepository(rdb)
	authService := services.NewAuthService(userRepo, session)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)
	r.POST("/verify", authHandler.VerifyToken)

	r.Run(":9085")
}