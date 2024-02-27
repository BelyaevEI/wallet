package app

import (
	"net"
	"net/http"

	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/logger"
	"github.com/BelyaevEI/wallet/internal/server/route"
	"github.com/BelyaevEI/wallet/internal/server/service"
)

// Application struct
type application struct {
	Server  *http.Server     // the server that processes requests for funds transfer
	Service service.Servicer // service for processing request
	// Beaver
}

func NewApp() (application, error) {

	// Create new connect to logger
	log, err := logger.NewLogger()
	if err != nil {
		return application{}, err
	}

	// Reading config file
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		log.Log.Error("read config file is fail: ", err)
		return application{}, err
	}

	// Init server service
	service, err := service.NewService(log, cfg)
	if err != nil {
		return application{}, err
	}

	// Init new router
	route, err := route.NewRouter(log, cfg, service)
	if err != nil {
		return application{}, err
	}

	// Init server
	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: route,
	}

	return application{Server: server}, nil
}
