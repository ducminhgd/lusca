package httpapi

import (
	"net/http"
	"strconv"

	"github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/internal/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Collection interface{}

type collectionRouter struct {
	repo *repositories.CollectionRepo
}

func NewCollectionRouter(db *db.GORMManager) *collectionRouter {
	return &collectionRouter{
		repo: repositories.NewCollectionRepo(db),
	}
}

func CollectionRouter(db *db.GORMManager) http.Handler {
	r := NewCollectionRouter(db)

	cr := chi.NewRouter()
	cr.Get("/by-id/{id}", r.GetByID)
	cr.Get("/by-name/{name}", r.GetByName)
	cr.Get("/", r.GetList)
	return cr
}

func (rt *collectionRouter) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	m, err := rt.repo.GetByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		NewNotFoundErrorResponse().Write(w)
		return
	}
	NewDataResponse(m, "Success").Write(w)
}

func (rt *collectionRouter) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	m, err := rt.repo.GetByName(r.Context(), name)
	if err != nil {
		NewNotFoundErrorResponse().Write(w)
		return
	}
	NewDataResponse(m, "Success").Write(w)
}

func (rt *collectionRouter) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := repositories.CollectionQuery{
		Name_Like: r.URL.Query().Get("name"),
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil || limit > 1000 {
		limit = 10
	}
	qOpts := repositories.QueryOptions{
		Limit:  int(limit),
		Offset: int((page - 1) * limit),
	}

	ls, count, err := rt.repo.GetList(ctx, q, qOpts)
	if err != nil {
		NewUnknownErrorResponse().Write(w)
		return
	}

	lb := ListBody{
		Total:    count,
		Records:  ls,
		Page:     int(page),
		PageSize: int(limit),
	}
	NewDataResponse(lb, "Success").Write(w)
}
