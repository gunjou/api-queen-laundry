package customer

import (
	"context"
	"fmt"
	"queen-laundry/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CREATE CUSTOMER
func (r *Repository) CreateCustomer(ctx context.Context, nama, noHp, alamat string) error {
	now := utils.GetNowWITA()

	query := `
		INSERT INTO customers (nama, no_hp, alamat, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query, nama, noHp, alamat, now, now)
	return err
}

// GET ALL CUSTOMERS
func (r *Repository) GetCustomers(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT id_customer, nama, no_hp, alamat
		FROM customers
		WHERE is_active = 1
		ORDER BY id_customer DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id int
		var nama, noHp, alamat string

		err := rows.Scan(&id, &nama, &noHp, &alamat)
		if err != nil {
			return nil, err
		}

		item := map[string]interface{}{
			"id_customer": id,
			"nama":        nama,
			"no_hp":       noHp,
			"alamat":      alamat,
		}

		result = append(result, item)
	}

	return result, nil
}


// GET CUSTOMER BY ID
func (r *Repository) GetCustomerByID(ctx context.Context, id int) (map[string]interface{}, error) {
	query := `
		SELECT id_customer, nama, no_hp, alamat
		FROM customers
		WHERE id_customer = $1 AND is_active = 1
	`

	row := r.db.QueryRow(ctx, query, id)

	var custID int
	var nama, noHp, alamat string

	err := row.Scan(&custID, &nama, &noHp, &alamat)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id_customer": custID,
		"nama":        nama,
		"no_hp":       noHp,
		"alamat":      alamat,
	}

	return result, nil
}

// UPDATE CUSTOMER BY ID
func (r *Repository) UpdateCustomer(ctx context.Context, id int, nama, noHp, alamat string) error {
	now := utils.GetNowWITA()

	query := `
		UPDATE customers
		SET nama = $1,
			no_hp = $2,
			alamat = $3,
			updated_at = $4
		WHERE id_customer = $5 AND is_active = 1
	`

	cmd, err := r.db.Exec(ctx, query, nama, noHp, alamat, now, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// DELETE CUSTOMER BY ID (SOFT DELETE)
func (r *Repository) DeleteCustomer(ctx context.Context, id int) error {
	now := utils.GetNowWITA()

	query := `
		UPDATE customers
		SET is_active = 0,
			updated_at = $1
		WHERE id_customer = $2
	`

	cmd, err := r.db.Exec(ctx, query, now, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}