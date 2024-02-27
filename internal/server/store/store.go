package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BelyaevEI/wallet/internal/config"
)

type Storer interface {
	CheckExists(ctx context.Context, id uint32) (bool, error)
	GetBalanceByID(ctx context.Context, id uint32) (int, error)
	CloseConnection2DB()
}

// Database layer
type Store struct {
	db *sql.DB
}

func Connect(cfg config.Config) (Storer, error) {

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	// Create table for wallet
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wallet
					(wallet_id bigint NOT NULL,
				    amount DECIMAL(10, 2) NOT NULL)`)
	if err != nil {
		return nil, err
	}

	return Store{db: db}, nil
}

func (store Store) CheckExists(ctx context.Context, id uint32) (bool, error) {

	var idEx uint32

	row := store.db.QueryRowContext(ctx, "SELECT wallet_id FROM wallet WHERE wallet_id = $1", id)
	if err := row.Scan(&idEx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (store Store) GetBalanceByID(ctx context.Context, id uint32) (int, error) {

	var balance int

	row := store.db.QueryRowContext(ctx, "SELECT wallet_id, amount FROM wallet WHERE wallet_id = $1", id)
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil

}

func (store Store) CloseConnection2DB() {
	store.db.Close()
}
