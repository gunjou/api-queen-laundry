package report

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// ================= WEEKLY REVENUE =================
// menampilkan pendapatan per hari dalam minggu berjalan
func (r *Repository) GetWeeklyRevenue(ctx context.Context) ([]map[string]interface{}, error) {

	query := `
		SELECT
			TO_CHAR(tanggal_bayar, 'Day') as label,
			COALESCE(SUM(jumlah), 0) as total
		FROM payments
		WHERE is_active = 1
		AND DATE(tanggal_bayar) >= DATE_TRUNC('week', CURRENT_DATE)
		AND DATE(tanggal_bayar) < DATE_TRUNC('week', CURRENT_DATE) + INTERVAL '7 day'
		GROUP BY DATE(tanggal_bayar), label
		ORDER BY DATE(tanggal_bayar)
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {

		var label string
		var total float64

		err := rows.Scan(&label, &total)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"label": label,
			"total": total,
		})
	}

	return result, nil
}

// ================= MONTHLY REVENUE =================
// menampilkan pendapatan per minggu dalam bulan tertentu
func (r *Repository) GetMonthlyRevenue(
	ctx context.Context,
	month string,
	year string,
) ([]map[string]interface{}, error) {

	query := `
		SELECT
			CONCAT('Minggu ', CEIL(EXTRACT(DAY FROM tanggal_bayar) / 7.0)) as label,
			COALESCE(SUM(jumlah), 0) as total
		FROM payments
		WHERE is_active = 1
		AND EXTRACT(MONTH FROM tanggal_bayar) = $1
		AND EXTRACT(YEAR FROM tanggal_bayar) = $2
		GROUP BY CEIL(EXTRACT(DAY FROM tanggal_bayar) / 7.0)
		ORDER BY CEIL(EXTRACT(DAY FROM tanggal_bayar) / 7.0)
	`

	rows, err := r.db.Query(ctx, query, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {

		var label string
		var total float64

		err := rows.Scan(&label, &total)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"label": label,
			"total": total,
		})
	}

	return result, nil
}

// ================= YEARLY REVENUE =================
// menampilkan pendapatan per bulan dalam tahun tertentu
func (r *Repository) GetYearlyRevenue(
	ctx context.Context,
	year string,
) ([]map[string]interface{}, error) {

	query := `
		SELECT
			TO_CHAR(tanggal_bayar, 'Month') as label,
			COALESCE(SUM(jumlah), 0) as total
		FROM payments
		WHERE is_active = 1
		AND EXTRACT(YEAR FROM tanggal_bayar) = $1
		GROUP BY EXTRACT(MONTH FROM tanggal_bayar), label
		ORDER BY EXTRACT(MONTH FROM tanggal_bayar)
	`

	rows, err := r.db.Query(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {

		var label string
		var total float64

		err := rows.Scan(&label, &total)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"label": label,
			"total": total,
		})
	}

	return result, nil
}