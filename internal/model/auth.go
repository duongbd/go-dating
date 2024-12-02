package model

import "github.com/golang-jwt/jwt"

// Claims is a struct that will be encoded to a JWT
type Claims struct {
	jwt.StandardClaims
	UserID uint
}
