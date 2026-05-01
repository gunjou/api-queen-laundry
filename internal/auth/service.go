package auth

import (
	"context"
	"fmt"
	"queen-laundry/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	hashed := user["password"].(string)

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// generate JWT
	token, err := utils.GenerateJWT(user["id_user"].(int), username)
	if err != nil {
		return "", err
	}

	return token, nil
}