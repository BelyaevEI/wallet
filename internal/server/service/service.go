// The service layer is designed to handle
// incoming requests to the server and business logic
package service

import (
	"net/http"

	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/logger"
	"github.com/BelyaevEI/wallet/internal/server/repository"
)

// Implementation check
var _ Servicer = &Service{}

type Servicer interface {
	Shutdown()
	GetBalance(writer http.ResponseWriter, request *http.Request)
	TransferFunds(writer http.ResponseWriter, request *http.Request)
}

type Service struct {
	log        *logger.Logger
	repository repository.Repositorer
}

func NewService(log *logger.Logger, cfg config.Config) (Servicer, error) {

	repository, err := repository.NewRepo(cfg)
	if err != nil {
		return &Service{}, err
	}

	return &Service{
		log:        log,
		repository: repository,
	}, nil
}

func (service Service) Shutdown() {
	service.repository.Shutdown()
}
