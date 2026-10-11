package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/Nainyxx/market-go/services/userservice/models"
	"github.com/Nainyxx/market-go/services/userservice/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(svc *service.Service) *UserHandler {
	return &UserHandler{service: svc}
}

type createUserRequest struct {
	ID      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
}

type updateUserRequest struct {
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	ImageUrl *string `json:"image_url"`
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, models.ErrAlreadyExist):
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
	case errors.Is(err, service.ErrInvalidUserID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, models.ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, models.ErrNotEnoughMoney):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		slog.Error("internal error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

// GET /users/:id
func (h *UserHandler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		user, err := h.service.GetUser(c.Request.Context(), user_id)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// POST /users
func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		userID, err := uuid.Parse(req.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		user, err := h.service.CreateUser(c.Request.Context(), userID, req.Name, req.Surname)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, user) // 201
	}
}

// PATCH /users/:id
func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var req updateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		user, err := h.service.UpdateUser(c.Request.Context(), &service.UpdateUserProfileInput{
			UserID:   userID,
			Name:     req.Name,
			Surname:  req.Surname,
			ImageUrl: req.ImageUrl,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
			writeError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// PATCH /users/addbalance/:id/:sum
func (h *UserHandler) AddBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		sum, err := strconv.ParseInt(c.Param("sum"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sum"})
			return
		}
		user, err := h.service.AddBalance(c.Request.Context(), userID, sum)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// Withdraw /users/withdraw/:id/:sum
func (h *UserHandler) Withdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		sum, err := strconv.ParseInt(c.Param("sum"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sum"})
			return
		}
		user, err := h.service.Withdraw(c.Request.Context(), userID, sum)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// AddRate /users/addrate/:id/:rate
func (h *UserHandler) AddRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		rate, err := strconv.ParseFloat(c.Param("rate"), 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rate"})
			return
		}
		user, err := h.service.AddRate(c.Request.Context(), userID, float32(rate))
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
