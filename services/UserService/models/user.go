package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Nainyxx/market-go/services/userservice/utils"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`

	Name     string `json:"name"`
	Surname  string `json:"surname"`
	ImageUrl string `json:"image_url"`
	// диапазон рейтинга от 0 до 5 включительно
	Rate       float32 `json:"rate"`
	RatesCount int32   `json:"rates_count"`

	// в копейках
	Balance int64 `json:"balance"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrUserCreate = errors.New("create user error")
	ErrUserChange = errors.New("change user error")
	ErrBalance    = errors.New("balance error")

	ErrSameName       = errors.New("name must be different")
	ErrSameSurname    = errors.New("surname must be different")
	ErrSameImageUrl   = errors.New("image_url must be different")
	ErrNotEnoughMoney = errors.New("not enough money on balance")
	ErrRateInterval   = errors.New("rate value must be from 0 to 5")
	ErrInvalidAmount  = errors.New("sum must be > 0")
	ErrNotFound       = errors.New("user not found")
)

func now() time.Time { return time.Now().UTC() }

func creationErr(msg error) error {
	return fmt.Errorf("%w: %w", ErrUserCreate, msg)
}
func changeErr(msg error) error {
	return fmt.Errorf("%w: %w", ErrUserChange, msg)
}
func balanceErr(msg error) error {
	return fmt.Errorf("%w: %w", ErrBalance, msg)
}

func CreateUser(id uuid.UUID, name, surname string) (*User, error) {
	err := utils.IsNameValid(name)
	if err != nil {
		return nil, creationErr(err)
	}
	err = utils.IsSurnameValid(surname)
	if err != nil {
		return nil, creationErr(err)
	}

	t := now()
	return &User{
		ID:        id,
		Name:      name,
		Surname:   surname,
		ImageUrl:  "",
		Rate:      0,
		Balance:   0,
		CreatedAt: t,
		UpdatedAt: t,
	}, nil
}

func (u *User) ChangeName(newName string) error {
	if newName == u.Name {
		return changeErr(ErrSameName)
	}
	err := utils.IsNameValid(newName)
	if err != nil {
		return changeErr(err)
	}

	u.Name = newName
	u.UpdatedAt = now()
	return nil
}

func (u *User) ChangeSurname(newSurname string) error {
	if newSurname == u.Surname {
		return changeErr(ErrSameSurname)
	}
	err := utils.IsSurnameValid(newSurname)
	if err != nil {
		return changeErr(err)
	}

	u.Surname = newSurname
	u.UpdatedAt = now()
	return nil
}

func (u *User) ChangeImageUrl(newImageUrl string) error {
	if newImageUrl == u.ImageUrl {
		return changeErr(ErrSameImageUrl)
	}
	err := utils.IsImageURLValid(newImageUrl)
	if err != nil {
		return changeErr(err)
	}

	u.ImageUrl = newImageUrl
	u.UpdatedAt = now()
	return nil
}

func (u *User) ChangeRate(newRate float32) error {
	if newRate < 0 || newRate > 5 {
		return changeErr(ErrRateInterval)
	}
	temp := u.Rate * float32(u.RatesCount)
	temp += newRate
	u.Rate = temp / float32(u.RatesCount+1)
	u.RatesCount++
	u.UpdatedAt = now()
	return nil
}

func (u *User) PlusToBalance(sum int64) error {
	if sum <= 0 {
		return balanceErr(ErrInvalidAmount)
	}
	u.Balance += sum
	u.UpdatedAt = now()
	return nil
}

func (u *User) MinusFromBalance(sum int64) error {
	if sum <= 0 {
		return balanceErr(ErrInvalidAmount)
	}
	if sum > u.Balance {
		return balanceErr(ErrNotEnoughMoney)
	}
	u.Balance -= sum
	u.UpdatedAt = now()
	return nil
}
