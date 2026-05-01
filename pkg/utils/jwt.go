package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("SECRET_LAUNDRY")

func GenerateJWT(userID int, username string) (string, error) {

	claims := jwt.MapClaims{
		"id_user":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SECRET_KEY)
}