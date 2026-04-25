package models

import (
	"github.com/golang-jwt/jwt/v5"
	
)

type UserAuth struct {
	Password    string `json:"password"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Id string `json:"id"`
}

type Profile struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Url         *string `json:"url,omitempty"`
	Role        string  `json:"role"`
}

type UserClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}
