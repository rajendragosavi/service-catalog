package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rajendragosavi/service-catalog/configs"
	"github.com/rajendragosavi/service-catalog/http/rest/handler"
	"github.com/rajendragosavi/service-catalog/pkg/db"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger *logrus.Logger
	router *mux.Router
	config configs.Config
}

func NewServer() (*Server, error) {
	cfg, err := configs.NewParsedConfig()
	if err != nil {
		return nil, err
	}

	database, err := db.Connect(db.ConfingDB{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	})
	if err != nil {
		return nil, err
	}

	log := NewLogger()
	router := mux.NewRouter()
	handler.Register(router, log, database)

	s := Server{
		logger: log,
		config: cfg,
		router: router,
	}
	return &s, nil
}

// run the server with graceful shutdown
func (s *Server) Run(ctx context.Context) error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.ServerPort),
		Handler: cors.Default().Handler(s.router),
	}
	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stopServer)

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("REST API listening on port %d", s.config.ServerPort)
		serverErrors <- server.ListenAndServe() // put any server errors on the error channel
	}(&wg)
	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received STOP signal")
		// shutdown the listern first
		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait() // wait for server shutdown
		s.logger.Info("server was shut down gracefully")
	}

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
