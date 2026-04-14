package service

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateService(ctx context.Context, nama string, tipe string, harga float64) error {
	return s.repo.CreateService(ctx, nama, tipe, harga)
}

func (s *Service) GetServices(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.GetServices(ctx)
}

func (s *Service) GetServiceByID(ctx context.Context, id int) (map[string]interface{}, error) {
	return s.repo.GetServiceByID(ctx, id)
}

func (s *Service) UpdateService(ctx context.Context, id int, nama, tipe string, harga float64) error {
	return s.repo.UpdateService(ctx, id, nama, tipe, harga)
}

func (s *Service) DeleteService(ctx context.Context, id int) error {
	return s.repo.DeleteService(ctx, id)
}