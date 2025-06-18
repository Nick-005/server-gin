// utils/token.go
package utils

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY_USER")) // Должен быть в конфиге!

type ClaimToRecover struct {
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	expiresAt time.Time `json:expires_at`
	jwt.RegisteredClaims
}

type ResetToken struct {
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
	Role      string    `json:"role"`
}

func GenerateResetToken(email, role string) (string, error) {

	claim := &ClaimToRecover{
		Role:      role,
		Email:     email,
		expiresAt: time.Now().Add(time.Minute * 30),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	result, err := token.SignedString(secretKey)
	if err != nil {
		return "error", err
	}
	result = strings.ReplaceAll(result, "+", "-")
	result = strings.ReplaceAll(result, "/", "_")
	return result, nil
}

func ValidateResetToken(tokenString string) (*ClaimToRecover, error) {

	claim := &ClaimToRecover{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {

		return nil, err
	}
	return claim, nil

}
