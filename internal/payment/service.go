package payment

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePayment(ctx context.Context, idOrder int, jumlah float64, metode string) error {

	validMetode := map[string]bool{
		"CASH":     true,
		"TRANSFER": true,
	}

	if !validMetode[metode] {
		return fmt.Errorf("invalid payment method")
	}

	return s.repo.CreatePayment(ctx, idOrder, jumlah, metode)
}

func (s *Service) GetPayments(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.GetPayments(ctx)
}

func (s *Service) GetPaymentByID(ctx context.Context, id int) (map[string]interface{}, error) {
	return s.repo.GetPaymentByID(ctx, id)
}