package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func RootRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", Health)
	return r
}
