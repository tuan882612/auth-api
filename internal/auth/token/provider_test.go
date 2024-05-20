package token

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const testSecert = "Secret"

func Test_Token_Generation(t *testing.T) {
	prov := NewProvider(testSecert)
	token, err := prov.GenerateToken(uuid.Nil, "some email")
	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("Token is empty")
	}

	t.Log(token)
}

func Test_Token_Decode(t *testing.T) {
	prov := NewProvider(testSecert)
	token, err := prov.GenerateToken(uuid.Nil, "some email")
	if err != nil {
		t.Error(err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecert), nil
	})

	if err != nil {
		t.Error(err)
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		t.Error("Failed to parse claims")
	}

	if claims.UserID != uuid.Nil.String() {
		t.Errorf("UserID is not equal, got: %s", claims.UserID)
	}

	if claims.Email != "some email" {
		t.Errorf("Email is not equal, got: %s", claims.Email)
	}

	t.Log(claims)
}
