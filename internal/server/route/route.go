package route

import (
	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/logger"
	"github.com/BelyaevEI/wallet/internal/server/service"
	"github.com/go-chi/chi"
)

func NewRouter(log *logger.Logger, cfg config.Config, service service.Servicer) (*chi.Mux, error) {

	// New router
	route := chi.NewRouter()

	// Handlers
	route.Get("/api/v1/{walletid}", service.GetBalance)

	return route, nil
}
