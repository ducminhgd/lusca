package httpapi

import (
	"log/slog"

	"github.com/ducminhgd/lusca/internal/repositories"
	"github.com/go-chi/httplog/v2"
)

type featureAPI struct {
	repo   *repositories.FeatureRepo
	logger *httplog.Logger
}

func NewFeatureRouter(dependencies ...interface{}) *featureAPI {
	api := &featureAPI{}

	for _, d := range dependencies {
		switch rp := d.(type) {
		case *repositories.FeatureRepo:
			api.repo = rp
		case *httplog.Logger:
			api.logger = rp
		default:
			slog.Warn("does not support this repo", slog.Any("repo", rp))
		}
	}

	return api
}

func (api *featureAPI) UseRepo(rp interface{}) *featureAPI {
	switch repo := rp.(type) {
	case *repositories.FeatureRepo:
		api.repo = repo
	default:
		slog.Warn("does not support this repo", slog.Any("repo", repo))
	}
	return api
}

func (api *featureAPI) UseLogger(logger *httplog.Logger) *featureAPI {
	api.logger = logger
	return api
}
