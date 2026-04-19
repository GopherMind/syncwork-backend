package jwt

import (
	"log"
	"os"
	"time"

	"github.com/GopherMind/syncwork-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

func Createjwt(u models.UserAuth) (string, error) {
	jwtSecretKey := []byte(os.Getenv("SECRET_KEY_JWT"))

	if len(jwtSecretKey) == 0 {
		log.Printf("WARNING: SECRET_KEY_JWT is empty or not set!")
	}

	claims := models.UserClaims{
		Id: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(720 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	return tokenString, nil
}