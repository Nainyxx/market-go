package repo

import (
	"context"

	"github.com/Nainyxx/market-go/services/userservice/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, newUser *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
