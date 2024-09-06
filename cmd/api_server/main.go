package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ducminhgd/lusca/config"
	"github.com/ducminhgd/lusca/internal"
	"github.com/ducminhgd/lusca/internal/httpapi"
	"github.com/ducminhgd/lusca/internal/repositories"
	"github.com/ducminhgd/lusca/internal/require"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func server(cfg *config.Config) http.Handler {
	apiLogger := internal.InitAPILogger(cfg)
	dbManager, err := require.NewDatabaseManager(cfg.Databases)
	if err != nil {
		panic(err)
	}

	// Declare repositories
	collectionRepo := repositories.NewCollectionRepo(dbManager)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	r.Use(httplog.RequestLogger(apiLogger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Mount("/", httpapi.RootRouter())

	// Collection API
	collectionAPI := httpapi.NewCollectionRouter(apiLogger, collectionRepo)
	r.Mount("/collections", collectionAPI.Router())
	return r
}

func main() {
	var (
		cfg = config.Load()
	)
	apiServer := server(&cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler: apiServer,
	}
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
