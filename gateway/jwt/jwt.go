package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	AccountId string `json:"account_id"`
	jwt.RegisteredClaims
}

func VerifyAndGetId(token string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", status.Error(codes.Unauthenticated, "Unauthorized request")
		}
		return "", status.Error(codes.InvalidArgument, err.Error())
	}
	if !tkn.Valid {
		return "", status.Error(codes.Unauthenticated, "Unauthorized request")
	}
	return claims.AccountId, nil
}
