package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (map[string]interface{}, error) {
	query := `
		SELECT id_user, username, password
		FROM users
		WHERE username = $1 AND is_active = 1
	`

	row := r.db.QueryRow(ctx, query, username)

	var id int
	var uname, password string

	err := row.Scan(&id, &uname, &password)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id_user":  id,
		"username": uname,
		"password": password,
	}, nil
}