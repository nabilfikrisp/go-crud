package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nabilfikrisp/go-crud/config"
	"github.com/nabilfikrisp/go-crud/internal/controller/restapi"
	"github.com/nabilfikrisp/go-crud/internal/repo/persistent"
	"github.com/nabilfikrisp/go-crud/internal/usecase/contact"
	"github.com/nabilfikrisp/go-crud/pkg/httpserver"
	"github.com/nabilfikrisp/go-crud/pkg/logger"
	"github.com/nabilfikrisp/go-crud/pkg/postgres"
)

type useCases struct {
	contact *contact.UseCase
}

type servers struct {
	http *httpserver.Server
}

func initUseCases(pg *postgres.Postgres) useCases {
	contactRepo := persistent.NewContactPGRepo(pg)

	return useCases{
		contact: contact.New(contactRepo),
	}
}

func initServer(cfg *config.Config, uc useCases, l logger.Interface) servers {
	httpserver := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpserver.Engine, cfg, uc.contact, l)

	return servers{
		http: httpserver,
	}
}

func (s *servers) startServers() {
	s.http.Start()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - shutdownServers - httpServer.Shutdown: %w", err))
	}
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	uc := initUseCases(pg)
	s := initServer(cfg, uc, l)
	s.startServers()
	s.waitForShutdown(l)
}
