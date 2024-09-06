package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ducminhgd/lusca/internal/models"
	"github.com/ducminhgd/lusca/internal/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
)

type CollectionRequestBody struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Values      []string `json:"values"`
}

type collectionAPI struct {
	repo   *repositories.CollectionRepo
	logger *httplog.Logger
}

// NewCollectionRouter creates a new collectionAPI with the given dependencies.
// The dependencies can be one or more of the following:
// - *repositories.CollectionRepo: a repository for collections
// - *httplog.Logger: a logger for the API
// All other dependencies are ignored.
func NewCollectionRouter(dependencies ...interface{}) *collectionAPI {
	api := &collectionAPI{}
	for _, d := range dependencies {
		switch rp := d.(type) {
		case *repositories.CollectionRepo:
			api.repo = rp
		case *httplog.Logger:
			api.logger = rp
		default:
			slog.Warn("does not support this repo", slog.Any("repo", rp))
		}
	}
	return api
}

// UseRepo injects a repository into the collectionAPI. If the repository is not
// a *repositories.CollectionRepo, it will be ignored and a warning will be
// logged.
func (api *collectionAPI) UseRepo(rp interface{}) *collectionAPI {
	switch repo := rp.(type) {
	case *repositories.CollectionRepo:
		api.repo = repo
	default:
		slog.Warn("does not support this repo", slog.Any("repo", repo))
	}
	return api
}

// UseLogger injects a logger into the collectionAPI. If the logger is not
// a *httplog.Logger, it will be ignored and a warning will be logged.
func (api *collectionAPI) UseLogger(logger *httplog.Logger) *collectionAPI {
	api.logger = logger
	return api
}

// Router returns a chi.Router which can be used to handle requests to the
// collection API endpoints.
func (api *collectionAPI) Router() http.Handler {
	cr := chi.NewRouter()
	cr.Get("/by-id/{id}", api.GetByID)
	cr.Get("/by-name/{name}", api.GetByName)
	cr.Get("/", api.GetList)
	cr.Post("/", api.Create)
	return cr
}

// GetByID is an HTTP handler which returns a collection by the given id.
//
// The id is a path parameter passed in the URL, and it must be a valid UUID.
// If the id is invalid, or the collection does not exist in the database, a
// 404 error is returned. Otherwise, the collection is returned as JSON in the
// response body.
func (api *collectionAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	m, err := api.repo.GetByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		NewNotFoundErrorResponse().Write(w)
		return
	}
	NewDataResponse(m, "Success").Write(w)
}

// GetByName is an HTTP handler which returns a collection by the given name.
//
// The name is a path parameter passed in the URL.
// If the name is invalid, or the collection does not exist in the database, a
// 404 error is returned. Otherwise, the collection is returned as JSON in the
// response body.
func (api *collectionAPI) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	m, err := api.repo.GetByName(r.Context(), name)
	if err != nil {
		NewNotFoundErrorResponse().Write(w)
		return
	}
	NewDataResponse(m, "Success").Write(w)
}

// GetList is an HTTP handler which returns a list of collections.
//
// The list of collections is queried based on the following parameters:
//   - name: a string that the collection's name must contain.
//   - page: the page number of the list to query. If not provided, it will default to 1.
//   - page_size: the number of items per page. If not provided, it will default to 10, and
//     will not exceed 1000.
//
// The response is a JSON object that contains the following fields:
// - total: the total number of collections that match the query.
// - records: a list of collections that match the query.
// - page: the page number of the list returned.
// - page_size: the number of items in the list returned.
func (api *collectionAPI) GetList(w http.ResponseWriter, r *http.Request) {
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

	ls, count, err := api.repo.GetList(ctx, q, qOpts)
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

func (api *collectionAPI) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CollectionRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		NewBadRequestResponse().Write(w)
		api.logger.Error("invalid JSON format", slog.Any("error", err))
		return
	}
	m := &models.Collection{
		Name:        req.Name,
		Description: req.Description,
	}
	for _, v := range req.Values {
		m.Details = append(m.Details, models.CollectionDetail{
			Value: v,
		})
	}
	affectedRows, err := api.repo.Create(ctx, m)
	api.logger.Info("create collection", slog.Any("affectedRows", affectedRows))
	if err != nil || affectedRows < 1 {
		NewUnknownErrorResponse().Write(w)
		api.logger.Error("failed to create collection", slog.Any("error", err))
		return
	}
	NewDataResponse(m, "Success").Write(w)
}
