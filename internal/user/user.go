package user

import (
	"io"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	LastLogin time.Time `json:"last_login"`
	Created   time.Time `json:"created"`
}

func (u *UserReq) DecodeValidate(r io.ReadCloser) error {
	if err := jsoniter.NewDecoder(r).Decode(u); err != nil {
		return err
	}

	if err := validator.New().Struct(u); err != nil {
		return err
	}

	return nil
}

func NewUser(email string, password []byte) *User {
	return &User{
		UserID:    uuid.New(),
		Email:     email,
		Password:  password,
		LastLogin: time.Now(),
		Created:   time.Now(),
	}
}

func (u *User) ValidatePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		log.Warn().Msgf("%s: invalid password attempt", u.UserID)
		return false
	}

	return true
}
