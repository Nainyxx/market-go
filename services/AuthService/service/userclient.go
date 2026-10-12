package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidProfile = errors.New("invalid name or surname")

type UserClient interface {
	CreateProfile(ctx context.Context, id uuid.UUID, name, surname string) error
}

type httpUserClient struct {
	baseURL string
	http    *http.Client
}

func NewHTTPUserClient(baseURL string) UserClient {
	return &httpUserClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 3 * time.Second},
	}
}

// POST {baseURL}/users;
func (c *httpUserClient) CreateProfile(ctx context.Context, id uuid.UUID, name, surname string) error {
	body, err := json.Marshal(map[string]string{
		"id":      id.String(),
		"name":    name,
		"surname": surname,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("user service unreachable: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return ErrInvalidProfile
	default:
		return fmt.Errorf("user service returned status %d", resp.StatusCode)
	}
}
