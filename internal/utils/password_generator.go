package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%^&*()-_=+,.?/:;{}[]~"
)

func HassPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func GeneratePassword(length int, useUpper, useDigits, useSpecial bool) (string, error) {
	var charPool string
	charPool += lowerChars

	if useUpper {
		charPool += upperChars
	}
	if useDigits {
		charPool += digitChars
	}
	if useSpecial {
		charPool += specialChars
	}

	if len(charPool) == 0 {
		return "", fmt.Errorf("no character pool defined")
	}

	var password strings.Builder
	poolLength := big.NewInt(int64(len(charPool)))

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, poolLength)
		if err != nil {
			return "", err
		}
		password.WriteByte(charPool[index.Int64()])
	}

	return password.String(), nil
}

// GenerateMediumPassword - генерирует пароль средней сложности (8-12 символов, буквы + цифры)
func GenerateMediumPassword() (string, error) {
	return GeneratePassword(10, true, true, false)
}
