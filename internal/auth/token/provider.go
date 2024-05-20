package token

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Provider struct {
	SECRET     []byte
	claimsPool sync.Pool
}

func NewProvider(secert string) *Provider {
	return &Provider{
		SECRET: []byte(secert),
		claimsPool: sync.Pool{
			New: func() interface{} {
				return &Claims{}
			},
		},
	}
}

func (p *Provider) GenerateToken(userID uuid.UUID, email string) (string, error) {
	defer p.claimsPool.Put(p.claimsPool.Get())
	claims := p.claimsPool.Get().(*Claims)

	claims.UserID = userID.String()
	claims.Email = email
	claims.Issuer = "authapi"
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(p.SECRET)
	if err != nil {
		log.Error().Err(err).Send()
		return "", err
	}

	return token, nil
}
