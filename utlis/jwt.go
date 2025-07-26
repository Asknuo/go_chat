package utlis

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(ID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": ID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("BMvfawCaCjlDOzLAoYDUxLZWGIzerY53VeIm03Fy6uE="))
}
