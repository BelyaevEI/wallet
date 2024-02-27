package repository

import (
	"context"

	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/server/store"
)

type Repositorer interface {
	CheckExists(ctx context.Context, id uint32) (bool, error)
	GetBalanceByID(ctx context.Context, id uint32) (int, error)
	Shutdown()
}

// Repository layer
type Repository struct {
	Store store.Storer
}

// Create new repository for service
func NewRepo(cfg config.Config) (Repositorer, error) {

	store, err := store.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return Repository{Store: store}, nil
}

// Check exists wallet by id
func (repository Repository) CheckExists(ctx context.Context, id uint32) (bool, error) {
	return repository.Store.CheckExists(ctx, id)
}

// Getting  balance by id
func (repository Repository) GetBalanceByID(ctx context.Context, id uint32) (int, error) {
	return repository.Store.GetBalanceByID(ctx, id)
}

// Closing open connection
func (repository Repository) Shutdown() {
	repository.Store.CloseConnection2DB()
}
