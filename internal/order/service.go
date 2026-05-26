package order

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

func (s *Service) CreateOrder(
	ctx context.Context,
	idCustomer, idService int,
	berat float64,
	catatan string,
	estimasiSelesai string,
	ongkir float64,
	metode *string,
	langsungBayar bool,
) (string, error) {

	return s.repo.CreateOrder(
		ctx,
		idCustomer,
		idService,
		berat,
		catatan,
		estimasiSelesai,
		ongkir,
		metode,
		langsungBayar,
	)
}

func (s *Service) GetOrders(
	ctx context.Context,
	page int,
	limit int,
) ([]map[string]interface{}, int, error) {

	return s.repo.GetOrders(ctx, page, limit)
}

func (s *Service) GetOrderByID(ctx context.Context, id int) (map[string]interface{}, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *Service) UpdateOrder(ctx context.Context, id int, berat float64, catatan string) error {
	return s.repo.UpdateOrder(ctx, id, berat, catatan)
}

func (s *Service) DeleteOrder(ctx context.Context, id int) error {
	return s.repo.DeleteOrder(ctx, id)
}

func (s *Service) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	// validasi status (simple validation)
	validStatus := map[string]bool{
		"DITERIMA": true,
		"DIPROSES": true,
		"SELESAI":  true,
		"DIAMBIL":  true,
	}

	if !validStatus[status] {
		return fmt.Errorf("invalid status")
	}

	return s.repo.UpdateOrderStatus(ctx, id, status)
}

func (s *Service) UpdateOrderPayment(ctx context.Context, id int, jumlah float64, metode string) error {

	validMetode := map[string]bool{
		"CASH":     true,
		"TRANSFER": true,
	}

	if !validMetode[metode] {
		return fmt.Errorf("invalid payment method")
	}

	return s.repo.UpdateOrderPayment(ctx, id, jumlah, metode)
}

func (s *Service) GetActiveOrders(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.GetActiveOrders(ctx)
}

func (s *Service) GetCompletedOrders(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.GetCompletedOrders(ctx)
}

func (s *Service) GetOrderSummary(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetOrderSummary(ctx)
}