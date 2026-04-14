package customer

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCustomer(ctx context.Context, nama, noHp, alamat string) error {
	return s.repo.CreateCustomer(ctx, nama, noHp, alamat)
}

func (s *Service) GetCustomers(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.GetCustomers(ctx)
}

func (s *Service) GetCustomerByID(ctx context.Context, id int) (map[string]interface{}, error) {
	return s.repo.GetCustomerByID(ctx, id)
}

func (s *Service) UpdateCustomer(ctx context.Context, id int, nama, noHp, alamat string) error {
	return s.repo.UpdateCustomer(ctx, id, nama, noHp, alamat)
}

func (s *Service) DeleteCustomer(ctx context.Context, id int) error {
	return s.repo.DeleteCustomer(ctx, id)
}