package app

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BelyaevEI/wallet/internal/beaver"
	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/logger"
	"github.com/BelyaevEI/wallet/internal/server/route"
	"github.com/BelyaevEI/wallet/internal/server/service"
)

// Application struct
type application struct {
	Server  *http.Server     // the server that processes requests for funds transfer
	Service service.Servicer // service for processing request
	Sigint  chan os.Signal   // channel for given signal for graceful shutdown
	Beaver  beaver.Beaverer  // the entity that performs the work of transferring funds
}

func NewApp() (application, error) {

	// Create new connect to logger
	log, err := logger.NewLogger()
	if err != nil {
		return application{}, err
	}

	// Reading config file
	cfg, err := config.LoadConfig("../")
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

	// Init beaver
	beaver, err := beaver.NewBeaver(log, cfg)
	if err != nil {
		return application{}, err
	}

	// Creating channel for graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	return application{Server: server, Sigint: sigint, Beaver: beaver}, nil
}
