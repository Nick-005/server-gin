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
	// jsonData, err := json.Marshal(token)
	// if err != nil {
	// 	return "", err
	// }

	// // Подписываем токен
	// mac := hmac.New(sha256.New, secretKey)
	// mac.Write(jsonData)
	// signature := mac.Sum(nil)

	// // Формат: данные.подпись
	// return fmt.Sprintf(
	// 	"%s.%s",
	// 	base64.URLEncoding.EncodeToString(jsonData),
	// 	base64.URLEncoding.EncodeToString(signature),
	// ), nil
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
	// parts := strings.Split(tokenStr, ".")
	// if len(parts) != 2 {
	// 	return nil, fmt.Errorf("invalid token format")
	// }

	// // Декодируем данные
	// jsonData, err := base64.URLEncoding.DecodeString(parts[0])
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid token encoding")
	// }

	// // Проверяем подпись
	// mac := hmac.New(sha256.New, secretKey)
	// mac.Write(jsonData)
	// expectedSignature := mac.Sum(nil)

	// actualSignature, err := base64.URLEncoding.DecodeString(parts[1])
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid signature encoding")
	// }

	// if !hmac.Equal(actualSignature, expectedSignature) {
	// 	return nil, fmt.Errorf("invalid token signature")
	// }

	// // Парсим данные
	// var token ResetToken
	// if err := json.Unmarshal(jsonData, &token); err != nil {
	// 	return nil, fmt.Errorf("invalid token data")
	// }

	// // Проверяем срок действия
	// if time.Now().After(token.ExpiresAt) {
	// 	return nil, fmt.Errorf("token expired")
	// }

	// return &token, nil
}
