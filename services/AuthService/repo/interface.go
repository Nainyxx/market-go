package repo

import (
	"context"

	"github.com/Nainyxx/market-go/services/authservice/models"
	"github.com/google/uuid"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
