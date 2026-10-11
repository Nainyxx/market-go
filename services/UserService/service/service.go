package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Nainyxx/market-go/services/userservice/models"
	"github.com/Nainyxx/market-go/services/userservice/repo"
)

var (
	ErrInvalidUserID = errors.New("user_id cannot be empty")
	ErrNilUpdate     = errors.New("update cannot be nil")
)

type Service struct {
	repo repo.UserRepo
}

type UpdateUserProfileInput struct {
	UserID   uuid.UUID
	Name     *string
	Surname  *string
	ImageUrl *string
}

func NewUserService(repo repo.UserRepo) *Service {
	return &Service{repo: repo}
}

func formatErr(err error) error {
	return fmt.Errorf("UserService: %w", err)
}

// GET
func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, formatErr(err)
	}

	return user, nil
}

// POST
func (s *Service) CreateUser(ctx context.Context, userID uuid.UUID, name, surname string) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}
	newUser, err := models.CreateUser(userID, name, surname)
	if err != nil {
		return nil, formatErr(err)
	}
	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, formatErr(err)
	}
	return newUser, nil
}

// PATCH
func (s *Service) UpdateUser(ctx context.Context, update *UpdateUserProfileInput) (*models.User, error) {
	if update == nil {
		return nil, formatErr(ErrNilUpdate)
	}
	if update.UserID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}

	user, err := s.repo.GetUser(ctx, update.UserID)
	if err != nil {
		return nil, formatErr(err)
	}

	changed := false
	if update.Name != nil && *update.Name != user.Name {
		err = user.ChangeName(*update.Name)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.Surname != nil && *update.Surname != user.Surname {
		err = user.ChangeSurname(*update.Surname)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.ImageUrl != nil && *update.ImageUrl != user.ImageUrl {
		err = user.ChangeImageUrl(*update.ImageUrl)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if !changed {
		return user, nil
	}

	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, formatErr(err)
	}
	return user, nil
}

func (s *Service) AddBalance(ctx context.Context, userID uuid.UUID, sum int64) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, formatErr(err)
	}
	err = user.PlusToBalance(sum)
	if err != nil {
		return nil, formatErr(err)
	}
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, formatErr(err)
	}
	return user, nil
}

func (s *Service) Withdraw(ctx context.Context, userID uuid.UUID, sum int64) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, formatErr(err)
	}
	err = user.MinusFromBalance(sum)
	if err != nil {
		return nil, formatErr(err)
	}
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, formatErr(err)
	}
	return user, nil
}

func (s *Service) AddRate(ctx context.Context, userID uuid.UUID, rate float32) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, formatErr(err)
	}
	err = user.ChangeRate(rate)
	if err != nil {
		return nil, formatErr(err)
	}
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, formatErr(err)
	}
	return user, nil
}

// DELETE
func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return formatErr(ErrInvalidUserID)
	}
	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return formatErr(err)
	}
	return nil
}
