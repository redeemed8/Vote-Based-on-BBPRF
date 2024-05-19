package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserClaim UserClaim
	jwt.StandardClaims
}

type UserClaim struct {
	Id    uint32
	Phone string
	UUID  string
}
