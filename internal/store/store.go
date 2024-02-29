// the storage layer processes all database
// requests and performs all operations inherent
// in the business logic with the data
package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Implementation check
var _ Storer = Store{}

type Storer interface {
	CheckExists(ctx context.Context, id uint32) (bool, error)
	GetBalanceByID(ctx context.Context, id uint32) (int, error)
	TransferFunds(ctx context.Context, mes models.Transfer) error
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

	var balance float64

	row := store.db.QueryRowContext(ctx, "SELECT amount FROM wallet WHERE wallet_id = $1", id)
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}
	return int(balance), nil

}

func (store Store) CloseConnection2DB() {
	store.db.Close()
}

func (store Store) TransferFunds(ctx context.Context, mes models.Transfer) error {

	// Begin transaction for transfer funds between wallet
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	// Block table
	_, err = tx.Exec("LOCK TABLE wallet IN ACCESS EXCLUSIVE MODE")
	if err != nil {
		return err
	}

	// Subtraction amount from wallet
	_, err = tx.ExecContext(ctx, "UPDATE wallet SET amount = amount - $1  WHERE wallet_id = $2", mes.Amount, mes.WalletIDFrom)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add amount to wallet
	_, err = tx.ExecContext(ctx, "UPDATE wallet SET amount = amount + $1  WHERE wallet_id = $2", mes.Amount, mes.WalletIDTo)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil

}
