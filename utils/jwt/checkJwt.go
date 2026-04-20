package jwt

import (
	"os"

	"github.com/GopherMind/syncwork-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

var secret_key = os.Getenv("SECRET_KEY_JWT")

func CheckJwt(tokenJwt string) (*models.UserClaims, error) {
	secretKey := os.Getenv("SECRET_KEY_JWT")
	if secretKey == "" {
		return nil, errors.New("SECRET_KEY_JWT is empty")
	}
	println("SECRET_KEY_JWT length:", len(secretKey))

	token, err := jwt.ParseWithClaims(tokenJwt, &models.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

