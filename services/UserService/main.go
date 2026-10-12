package main

import (
	"log"

	"github.com/Nainyxx/market-go/services/userservice/conf"
	"github.com/Nainyxx/market-go/services/userservice/handler"
	"github.com/Nainyxx/market-go/services/userservice/repo"
	"github.com/Nainyxx/market-go/services/userservice/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := conf.Load()

	db, err := repo.ConnectDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repo.NewPostgresRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	router.GET("/users/:id", userHandler.GetUser())
	router.POST("/users", userHandler.CreateUser())
	router.PATCH("/users/:id", userHandler.UpdateUser())
	router.DELETE("/users/:id", userHandler.DeleteUser())
	router.PATCH("/users/addbalance/:id/:sum", userHandler.AddBalance())
	router.PATCH("/users/withdraw/:id/:sum", userHandler.Withdraw())
	router.PATCH("/users/addrate/:id/:rate", userHandler.AddRate())

	log.Println("user service started on port", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}
