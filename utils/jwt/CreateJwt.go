package jwt

import (
	
	"os"
	"time"

	"github.com/GopherMind/syncwork-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY_JWT"))

func Createjwt(u models.UserAuth) (string, error) {

	claims := models.UserClaims{
		Name: u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}