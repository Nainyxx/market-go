package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Nainyxx/market-go/services/authservice/models"
	"github.com/Nainyxx/market-go/services/authservice/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrProfileCreate      = errors.New("failed to create profile")
)

type Service struct {
	repo      repo.AuthRepo
	users     UserClient
	jwtSecret []byte
	accessTTL time.Duration
}

func formatErr(msg error) error {
	return fmt.Errorf("AuthService: %w", msg)
}

func NewAuthService(r repo.AuthRepo, users UserClient, secret string, ttl time.Duration) *Service {
	return &Service{repo: r, users: users, jwtSecret: []byte(secret), accessTTL: ttl}
}

func normalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

func (s *Service) Register(ctx context.Context, phone, email, password, name, surname string) (*models.User, error) {
	user, err := models.CreateUser(strings.TrimSpace(phone), normalizeEmail(email), password)
	if err != nil {
		return nil, formatErr(err)
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, formatErr(err)
	}

	if err := s.users.CreateProfile(ctx, user.ID, name, surname); err != nil {
		// компенсация: откатываем пользователя в своей БД
		rbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
		defer cancel()
		if delErr := s.repo.DeleteUser(rbCtx, user.ID); delErr != nil {
			slog.Error("rollback failed", "user_id", user.ID, "err", delErr)
		}
		return nil, formatErr(fmt.Errorf("%w: %w", ErrProfileCreate, err))
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", formatErr(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	return s.issueToken(uuid.UUID(user.ID))
}

func (s *Service) issueToken(id uuid.UUID) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   id.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
