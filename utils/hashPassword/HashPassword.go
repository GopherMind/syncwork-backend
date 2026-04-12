package hashpassword

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("error to hash password")
		return "", err
	}
	log.Println(hashedPassword)
	return string(hashedPassword), nil
}
