package auth

import (
	"context"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"authapi/internal/auth/token"
	"authapi/internal/user"
)

type SFAService struct {
	userRepo  user.UserRepository
	tokenProv *token.Provider
}

func NewSFAService(repo user.UserRepository, prov *token.Provider) *SFAService {
	return &SFAService{
		userRepo:  repo,
		tokenProv: prov,
	}
}

func (a *SFAService) Login(ctx context.Context, email, password string) (*user.User, string, error) {
	usr, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if !usr.ValidatePassword(password) {
		return nil, "", user.ErrInvalidPassword
	}

	token, err := a.tokenProv.GenerateToken(usr.UserID, usr.Email)
	if err != nil {
		return nil, "", err
	}

	return usr, token, nil
}

func (a *SFAService) Register(ctx context.Context, email, password string) (*user.User, string, error) {
	hashedPsw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, "", err
	}

	newUser := user.NewUser(email, hashedPsw)

	tx, err := a.userRepo.NewTx(ctx)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(ctx)

	if err := a.userRepo.Save(ctx, tx, newUser); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Send()
		return nil, "", err
	}

	token, err := a.tokenProv.GenerateToken(newUser.UserID, newUser.Email)
	if err != nil {
		return nil, "", err
	}

	return newUser, token, err
}
