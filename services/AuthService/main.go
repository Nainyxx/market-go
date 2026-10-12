package main

import (
	"log"

	"github.com/Nainyxx/market-go/services/authservice/conf"
	"github.com/Nainyxx/market-go/services/authservice/handler"
	"github.com/Nainyxx/market-go/services/authservice/repo"
	"github.com/Nainyxx/market-go/services/authservice/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := conf.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := repo.ConnectDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authRepo := repo.NewPostgresRepo(db)
	userClient := service.NewHTTPUserClient(cfg.UserServiceURL)
	authService := service.NewAuthService(authRepo, userClient, cfg.JWTSecret, cfg.AccessTTL)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	router.POST("/auth/register", authHandler.Register())
	router.POST("/auth/login", authHandler.Login())

	log.Println("auth service started on port", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}
