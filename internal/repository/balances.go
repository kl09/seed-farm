package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalancesRepository struct {
	db *pgxpool.Pool
}

func NewBalancesRepository(
	db *pgxpool.Pool,
) *BalancesRepository {
	return &BalancesRepository{
		db: db,
	}
}

func (r *BalancesRepository) Exists(ctx context.Context, address string) (bool, error) {
	var exists int
	err := r.db.QueryRow(
		ctx,
		`SELECT 1 FROM wallets WHERE address = $1`,
		address,
	).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("query row: %w", err)
	}
	return true, nil
}
