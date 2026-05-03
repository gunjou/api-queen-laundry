package order

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

// ================= CREATE ORDER =================
// CREATE ORDER
func (r *Repository) CreateOrder(
	ctx context.Context,
	idCustomer, idService int,
	berat float64,
	catatan string,
	metode *string,
	langsungBayar bool,
) (string, error) {

	now := utils.GetNowWITA()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// ================= AMBIL HARGA =================
	var hargaService float64
	err = tx.QueryRow(ctx,
		`SELECT harga FROM services WHERE id_service = $1 AND is_active = 1`,
		idService,
	).Scan(&hargaService)
	if err != nil {
		return "", err
	}

	hargaFinal := berat * hargaService

	invoice := fmt.Sprintf("INV-%s", now.Format("060102150405"))

	// ================= SET PAYMENT STATUS =================
	paymentStatus := "BELUM_BAYAR"
	if langsungBayar {
		paymentStatus = "SUDAH_BAYAR"
	}

	// ================= INSERT ORDER =================
	var orderID int

	queryOrder := `
		INSERT INTO orders (
			kode_invoice,
			id_customer,
			id_service,
			berat,
			harga,
			harga_final,
			order_status,
			payment_status,
			catatan,
			tanggal_masuk,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,
			'DITERIMA',
			$7,
			$8,
			$9,$9,$9
		)
		RETURNING id_order, kode_invoice
	`

	err = tx.QueryRow(ctx, queryOrder,
		invoice,
		idCustomer,
		idService,
		berat,
		hargaService,
		hargaFinal,
		paymentStatus,
		catatan,
		now,
	).Scan(&orderID, &invoice)

	if err != nil {
		return "", err
	}

	// ================= INSERT PAYMENT (HANYA JIKA BAYAR) =================
	if langsungBayar {
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

		_, err = tx.Exec(ctx,
			queryPayment,
			orderID,
			hargaFinal,
			metode,
			now,
		)

		if err != nil {
			return "", err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return invoice, nil
}


// GET ALL ORDER
func (r *Repository) GetOrders(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT o.id_order, o.kode_invoice, o.berat, o.harga_final,
		       o.order_status, o.payment_status,
		       c.nama, s.nama
		FROM orders o
		LEFT JOIN customers c ON o.id_customer = c.id_customer
		JOIN services s ON o.id_service = s.id_service
		WHERE o.is_active = 1
		ORDER BY o.id_order DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id int
		var invoice, status, payStatus string
		var namaCustomer, namaService *string
		var berat, harga float64

		err := rows.Scan(&id, &invoice, &berat, &harga,
			&status, &payStatus, &namaCustomer, &namaService)
		if err != nil {
			return nil, err
		}

		item := map[string]interface{}{
			"id_order":       id,
			"kode_invoice":   invoice,
			"berat":          berat,
			"harga_final":    harga,
			"order_status":   status,
			"payment_status": payStatus,
			"customer":       namaCustomer,
			"service":        namaService,
		}

		result = append(result, item)
	}

	return result, nil
}


// GET ORDER BY ID
func (r *Repository) GetOrderByID(ctx context.Context, id int) (map[string]interface{}, error) {
	query := `
		SELECT id_order, kode_invoice, berat, harga_final, order_status, payment_status
		FROM orders
		WHERE id_order = $1 AND is_active = 1
	`

	row := r.db.QueryRow(ctx, query, id)

	var oid int
	var invoice, status, payStatus string
	var berat, harga float64

	err := row.Scan(&oid, &invoice, &berat, &harga, &status, &payStatus)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id_order":       oid,
		"kode_invoice":   invoice,
		"berat":          berat,
		"harga_final":    harga,
		"order_status":   status,
		"payment_status": payStatus,
	}, nil
}


// UPDATE ORDER BY ID
func (r *Repository) UpdateOrder(ctx context.Context, id int, berat float64, catatan string) error {
	now := utils.GetNowWITA()

	query := `
		UPDATE orders
		SET berat = $1,
		    catatan = $2,
		    updated_at = $3
		WHERE id_order = $4 AND is_active = 1
	`

	cmd, err := r.db.Exec(ctx, query, berat, catatan, now, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}


// DELETE ORDER BY ID (SOFT DELETE)
func (r *Repository) DeleteOrder(ctx context.Context, id int) error {
	now := utils.GetNowWITA()

	query := `
		UPDATE orders
		SET is_active = 0,
		    updated_at = $1
		WHERE id_order = $2
	`

	cmd, err := r.db.Exec(ctx, query, now, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}


// UPDATE ORDER STATUS BY ID
func (r *Repository) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	now := utils.GetNowWITA()
	
	query := `
	UPDATE orders
	SET order_status = $1,
	updated_at = $2
	WHERE id_order = $3 AND is_active = 1
	`
	
	cmd, err := r.db.Exec(ctx, query, status, now, id)
	if err != nil {
		return err
	}
	
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}
	
	return nil
}


// UPDATE ORDER PAYMENT BY ID
func (r *Repository) UpdateOrderPayment(ctx context.Context, id int, jumlah float64, metode string) error {
	now := utils.GetNowWITA()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// insert ke payments
	queryPayment := `
		INSERT INTO payments (id_order, jumlah, metode, tanggal_bayar, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4, $4)
	`

	_, err = tx.Exec(ctx, queryPayment, id, jumlah, metode, now)
	if err != nil {
		return err
	}

	// update order
	queryOrder := `
		UPDATE orders
		SET payment_status = 'SUDAH_BAYAR',
		    updated_at = $1
		WHERE id_order = $2 AND is_active = 1
	`

	cmd, err := tx.Exec(ctx, queryOrder, now, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return tx.Commit(ctx)
}



// GET ORDER BY STATUS = AKTIF (BELUM DIAMBIL)
func (r *Repository) GetActiveOrders(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT id_order, kode_invoice, order_status, payment_status
		FROM orders
		WHERE is_active = 1
		AND order_status != 'DIAMBIL'
		ORDER BY id_order DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id int
		var invoice, status, payment string

		err := rows.Scan(&id, &invoice, &status, &payment)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id_order":       id,
			"kode_invoice":   invoice,
			"order_status":   status,
			"payment_status": payment,
		})
	}

	return result, nil
}


func (r *Repository) GetCompletedOrders(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT id_order, kode_invoice, order_status, payment_status
		FROM orders
		WHERE is_active = 1
		AND order_status != 'DIAMBIL'
		ORDER BY id_order DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id int
		var invoice, status, payment string

		err := rows.Scan(&id, &invoice, &status, &payment)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id_order":       id,
			"kode_invoice":   invoice,
			"order_status":   status,
			"payment_status": payment,
		})
	}

	return result, nil
}


// GET ORDER BY SUMMARY
func (r *Repository) GetOrderSummary(ctx context.Context) (map[string]interface{}, error) {
	now := utils.GetNowWITA()

	query := `
		SELECT 
			COUNT(*) FILTER (WHERE DATE(tanggal_masuk) = DATE($1)) as total_today,
			COUNT(*) FILTER (WHERE order_status = 'DIAMBIL' AND DATE(tanggal_masuk) = DATE($1)) as completed_today,
			COUNT(*) FILTER (WHERE order_status != 'DIAMBIL') as active_orders,
			COALESCE(SUM(harga_final) FILTER (WHERE payment_status = 'SUDAH_BAYAR' AND DATE(tanggal_masuk) = DATE($1)), 0) as revenue_today
		FROM orders
		WHERE is_active = 1
	`

	row := r.db.QueryRow(ctx, query, now)

	var totalToday, completedToday, activeOrders int
	var revenueToday float64

	err := row.Scan(&totalToday, &completedToday, &activeOrders, &revenueToday)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_today":     totalToday,
		"completed_today": completedToday,
		"active_orders":   activeOrders,
		"revenue_today":   revenueToday,
	}, nil
}