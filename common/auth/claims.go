package auth

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}
