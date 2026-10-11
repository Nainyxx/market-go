package handler

import (
	"fmt"

	"github.com/Nainyxx/market-go/services/userservice/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service &service.Service
}

func formatErr(msg error) error {
	return fmt.Errorf("Handler: %w", msg)
}

func NewUserHandler() *UserHandler {
	serv := service.NewUserService()
	return &UserHandler{serv}
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {

}
