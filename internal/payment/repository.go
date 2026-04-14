package payment

import (
	"context"
	"fmt"
	"queen-laundry/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CREATE PAYMENT
func (r *Repository) CreatePayment(ctx context.Context, idOrder int, jumlah float64, metode string) error {
	now := utils.GetNowWITA()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// cek apakah sudah bayar
	var status string
	err = tx.QueryRow(ctx,
		`SELECT payment_status FROM orders WHERE id_order = $1 AND is_active = 1`,
		idOrder,
	).Scan(&status)

	if err != nil {
		return err
	}

	if status == "SUDAH_BAYAR" {
		return fmt.Errorf("order already paid")
	}

	// insert payment
	queryPayment := `
		INSERT INTO payments (
			id_order,
			jumlah,
			metode,
			tanggal_bayar,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$4,$4)
	`

	_, err = tx.Exec(ctx, queryPayment, idOrder, jumlah, metode, now)
	if err != nil {
		return err
	}

	// update order
	queryOrder := `
		UPDATE orders
		SET payment_status = 'SUDAH_BAYAR',
		    updated_at = $1
		WHERE id_order = $2
	`

	_, err = tx.Exec(ctx, queryOrder, now, idOrder)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}


// GET ALL PAYMENTS
func (r *Repository) GetPayments(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT p.id_payment, p.jumlah, p.metode, p.tanggal_bayar,
		       o.kode_invoice
		FROM payments p
		JOIN orders o ON p.id_order = o.id_order
		WHERE p.is_active = 1
		ORDER BY p.id_payment DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id int
		var jumlah float64
		var metode, invoice string
		var tanggal time.Time

		err := rows.Scan(&id, &jumlah, &metode, &tanggal, &invoice)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id_payment":    id,
			"jumlah":        jumlah,
			"metode":        metode,
			"tanggal_bayar": tanggal,
			"kode_invoice":  invoice,
		})
	}

	return result, nil
}


// GET PAYMENT BY ID
func (r *Repository) GetPaymentByID(ctx context.Context, id int) (map[string]interface{}, error) {
	query := `
		SELECT p.id_payment, p.jumlah, p.metode, p.tanggal_bayar,
		       o.kode_invoice
		FROM payments p
		JOIN orders o ON p.id_order = o.id_order
		WHERE p.id_payment = $1 AND p.is_active = 1
	`

	row := r.db.QueryRow(ctx, query, id)

	var pid int
	var jumlah float64
	var metode, invoice string
	var tanggal time.Time

	err := row.Scan(&pid, &jumlah, &metode, &tanggal, &invoice)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id_payment":    pid,
		"jumlah":        jumlah,
		"metode":        metode,
		"tanggal_bayar": tanggal,
		"kode_invoice":  invoice,
	}, nil
}