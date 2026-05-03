package report

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetRevenueReport(
	ctx context.Context,
	reportType string,
	month string,
	year string,
) ([]map[string]interface{}, error) {

	switch reportType {

	case "weekly":
		return s.repo.GetWeeklyRevenue(ctx)

	case "monthly":
		return s.repo.GetMonthlyRevenue(ctx, month, year)

	case "yearly":
		return s.repo.GetYearlyRevenue(ctx, year)

	default:
		return nil, fmt.Errorf("invalid report type")
	}
}