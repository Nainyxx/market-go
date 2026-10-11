package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Nainyxx/market-go/services/userservice/models"
	"github.com/lib/pq"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) UserRepo {
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

// GetUser
func (post *postgresRepo) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT
		user_id, name, surname, image_url, rate, rates_count, balance, created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	profile := &models.User{}
	err := post.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Surname,
		&profile.ImageUrl,
		&profile.Rate,
		&profile.RatesCount,
		&profile.Balance,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, formatErr(models.ErrNotFound)
	}
	if err != nil {
		return nil, formatErr(fmt.Errorf("get user: %w", err))
	}

	return profile, nil
}

func (post *postgresRepo) CreateUser(ctx context.Context, newUser *models.User) error {
	query := `
		INSERT INTO users (
		user_id, name, surname, image_url, rate, rates_count, balance, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := post.db.ExecContext(ctx, query,
		&newUser.ID,
		&newUser.Name,
		&newUser.Surname,
		&newUser.ImageUrl,
		&newUser.Rate,
		&newUser.RatesCount,
		&newUser.Balance,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return formatErr(fmt.Errorf("%w: %w", models.ErrAlreadyExist, err))
		}
		return formatErr(fmt.Errorf("failed to create user: %w", err))
	}

	return nil
}

func (post *postgresRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, surname = $2, image_url = $3, rate = $4, rates_count = $5, balance = $6, updated_at = $7
		WHERE user_id = $8
	`

	result, err := post.db.ExecContext(ctx, query,
		&user.Name,
		&user.Surname,
		&user.ImageUrl,
		&user.Rate,
		&user.RatesCount,
		&user.Balance,
		&user.UpdatedAt,
		&user.ID,
	)
	if err != nil {
		return formatErr(fmt.Errorf("failed to update user: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return formatErr(fmt.Errorf("failed to update user: %w", err))
	}
	if rowsAffected == 0 {
		return formatErr(models.ErrNotFound)
	}

	return nil
}

// DeleteUser
func (post *postgresRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE user_id = $1
	`

	result, err := post.db.ExecContext(ctx, query, userID)
	if err != nil {
		return formatErr(fmt.Errorf("failed to delete user: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return formatErr(fmt.Errorf("failed to delete user: %w", err))
	}
	if rowsAffected == 0 {
		return formatErr(models.ErrNotFound)
	}

	return nil
}
