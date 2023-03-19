package service

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	AccountId string `json:"account_id"`
	jwt.RegisteredClaims
}

func createToken(accountId string) (string, error) {
	claims := &Claims{
		AccountId: accountId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
