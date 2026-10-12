package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nainyxx/market-go/services/authservice/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) AuthRepo {
	return &postgresRepo{db: db}
}

func formatErr(msg error) error {
	return fmt.Errorf("Repo: %w", msg)
}

func ConnectDB(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	dbs := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dbs)
	if err != nil {
		return nil, formatErr(fmt.Errorf("error with connection: %w", err))
	}

	err = db.Ping()
	if err != nil {
		return nil, formatErr(fmt.Errorf("failed to ping db: %w", err))
	}

	return db, nil
}

// CreateUser
func (p *postgresRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
		user_id, phone, email, password_hash
		)
		VALUES ($1, $2, $3, $4)
	`
	_, err := p.db.ExecContext(ctx, query, user.ID, user.Phone, user.Email, user.PasswordHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return formatErr(fmt.Errorf("%w: %w", models.ErrAlreadyExist, err))
		}
		return formatErr(fmt.Errorf("failed to create user: %w", err))
	}

	return nil
}

// GetUserByEmail
func (p *postgresRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT user_id, phone, email, password_hash
		FROM users
		WHERE email = $1
	`
	var u models.User
	err := p.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Phone, &u.Email, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, formatErr(models.ErrNotFound)
		}
		return nil, formatErr(fmt.Errorf("failed to get user: %w", err))
	}

	return &u, nil
}

// GetUserByPhone
func (p *postgresRepo) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT user_id, phone, email, password_hash
		FROM users
		WHERE phone = $1
	`
	var u models.User
	err := p.db.QueryRowContext(ctx, query, phone).Scan(&u.ID, &u.Phone, &u.Email, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, formatErr(models.ErrNotFound)
		}
		return nil, formatErr(fmt.Errorf("failed to get user: %w", err))
	}

	return &u, nil
}

func (p *postgresRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE user_id = $1
	`
	_, err := p.db.ExecContext(ctx, query, userID)
	if err != nil {
		return formatErr(fmt.Errorf("failed to delete user: %w", err))
	}

	return nil
}
