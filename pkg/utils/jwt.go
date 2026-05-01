package utils

import (
	"queen-laundry/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateJWT(userID int, username string) (string, error) {
	cfg := config.LoadConfig()

	claims := jwt.MapClaims{
		"id_user":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(cfg.JwtSecret)
}