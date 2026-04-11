package models
import (
	"github.com/golang-jwt/jwt/v5"
)
type UserAuth struct {
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
	Role string `json:"role"`
	Url string `json:"url"`
	Description string `json:"description"`
	
}
type UserClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}