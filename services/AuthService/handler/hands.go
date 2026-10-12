package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Nainyxx/market-go/services/authservice/models"
	"github.com/Nainyxx/market-go/services/authservice/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.Service
}

func NewAuthHandler(svc *service.Service) *AuthHandler {
	return &AuthHandler{service: svc}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Surname  string `json:"surname" binding:"required,max=64"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func writeError(c *gin.Context, err error) {
	var createErr *models.CreateError
	switch {
	case errors.As(err, &createErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + createErr.Field})
	case errors.Is(err, models.ErrAlreadyExist):
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
	case errors.Is(err, service.ErrInvalidProfile):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name or surname"})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
	default:
		slog.Error("internal error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		user, err := h.service.Register(c.Request.Context(), req.Phone, req.Email, req.Password, req.Name, req.Surname)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": token})
	}
}
