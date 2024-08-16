package main

import (
	"net/http"

	"github.com/ducminhgd/lusca/config"
	"github.com/ducminhgd/lusca/internal/httpapi"
	"github.com/ducminhgd/lusca/internal/require"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func server() {
	var (
		cfg = config.Load()
	)
	dbManager, err := require.NewDatabaseManager(cfg.Databases)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Mount("/", httpapi.RootRouter())
	r.Mount("/collections", httpapi.CollectionRouter(dbManager))

	http.ListenAndServe(":8080", r)
}

func main() {
	server()
}
