package service

import (
	"context"
	"queen-laundry/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// INSERT SERVICE
func (r *Repository) CreateService(ctx context.Context, nama string, tipe string, harga float64) error {
	now := utils.GetNowWITA()

	query := `
		INSERT INTO services (nama, tipe, harga, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, nama, tipe, harga, now, now)
	return err
}

// GET ALL SERVICES
func (r *Repository) GetServices(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT id_service, nama, tipe, harga
		FROM services
		WHERE is_active = 1
		ORDER BY id_service DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id int
		var nama, tipe string
		var harga float64

		err := rows.Scan(&id, &nama, &tipe, &harga)
		if err != nil {
			return nil, err
		}

		item := map[string]interface{}{
			"id_service": id,
			"nama":       nama,
			"tipe":       tipe,
			"harga":      harga,
		}

		result = append(result, item)
	}

	return result, nil
}


// GET SERVICE BY ID
func (r *Repository) GetServiceByID(ctx context.Context, id int) (map[string]interface{}, error) {
	query := `
		SELECT id_service, nama, tipe, harga
		FROM services
		WHERE id_service = $1 AND is_active = 1
	`

	row := r.db.QueryRow(ctx, query, id)

	var svcID int
	var nama, tipe string
	var harga float64

	err := row.Scan(&svcID, &nama, &tipe, &harga)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id_service": svcID,
		"nama":       nama,
		"tipe":       tipe,
		"harga":      harga,
	}

	return result, nil
}

// UPDATE SERVICE
func (r *Repository) UpdateService(ctx context.Context, id int, nama, tipe string, harga float64) error {
	now := utils.GetNowWITA()

	query := `
		UPDATE services
		SET nama = $1,
			tipe = $2,
			harga = $3,
			updated_at = $5
		WHERE id_service = $4 AND is_active = 1
	`

	_, err := r.db.Exec(ctx, query, nama, tipe, harga, id, now)
	return err
}

// SOFT DELETE
func (r *Repository) DeleteService(ctx context.Context, id int) error {
	now := utils.GetNowWITA()
	
	query := `
		UPDATE services
		SET is_active = 0,
			updated_at = $2
		WHERE id_service = $1
	`

	_, err := r.db.Exec(ctx, query, id, now)
	return err
}