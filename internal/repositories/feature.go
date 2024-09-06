package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/internal"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type FeatureQuery struct {
	Name             string `json:"name" default:""`
	Name_Like        string `json:"name__like" default:""`
	Description      string `json:"description" default:""`
	Description_Like string `json:"description__like" default:""`
	Status           int    `json:"status" default:"0"`
	Status_In        []int  `json:"status__in"`
}

type FeatureRepo struct {
	db *db.GORMManager
}

func NewFeatureRepo(db *db.GORMManager) *FeatureRepo {
	return &FeatureRepo{
		db: db,
	}
}

func (r *FeatureRepo) GetList(ctx context.Context, conditions FeatureQuery, opts QueryOptions) ([]models.Feature, error) {
	var fs []models.Feature
	q := r.db.DB().WithContext(ctx).Model(&models.Feature{})
	if conditions.Name != "" {
		q = q.Where("name = @name", sql.Named("name", conditions.Name))
	}
	if conditions.Name_Like != "" {
		q = q.Where("name like @name__like", sql.Named("name__like", "%"+conditions.Name_Like+"%"))
	}
	if conditions.Description != "" {
		q = q.Where("description = @description", sql.Named("description", conditions.Description))
	}
	if conditions.Description_Like != "" {
		q = q.Where("description like @description__like", sql.Named("description__like", "%"+conditions.Description_Like+"%"))
	}
	if conditions.Status != 0 {
		q = q.Where("status = @status", sql.Named("status", conditions.Status))
	}
	if len(conditions.Status_In) > 0 {
		q = q.Where("status in (@status__in)", sql.Named("status__in", conditions.Status_In))
	}
	err := q.Order("created_at desc").
		Limit(opts.Limit).Offset(opts.Offset).Find(&fs).Error
	return fs, err
}

func (r *FeatureRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Feature, error) {
	var f models.Feature
	err := r.db.DB().WithContext(ctx).Model(&models.Feature{}).Where("id = ?", id).First(&f).Error
	return &f, err
}

func (r *FeatureRepo) GetByName(ctx context.Context, name string) (*models.Feature, error) {
	var f models.Feature
	err := r.db.DB().WithContext(ctx).Model(&models.Feature{}).Where("name = ?", name).First(&f).Error
	return &f, err
}

func (r *FeatureRepo) Create(ctx context.Context, m *models.Feature) (*models.Feature, error) {
	err := r.db.DB().WithContext(ctx).Create(&m).Error
	return m, err
}

func (r *FeatureRepo) Update(ctx context.Context, m *models.Feature) (*models.Feature, error) {
	err := r.db.DB().WithContext(ctx).Model(&models.Feature{}).Where("id = ?", m.ID).Updates(&m).Error
	return m, err
}

type FeatureEnableQuery struct {
	FeatureID       uuid.UUID `url:"feature_id" json:"feature_id" default:""`
	FeatureName     string    `url:"feature_name" json:"feature_name" default:""`
	Key             string    `url:"key" json:"key" default:""`
	Value           string    `url:"value" json:"value" default:""`
	Environment     string    `url:"environment" json:"environment" default:""`
	CollectionValue string    `url:"collection_value" json:"collection_value" default:""`
}

func (r *FeatureRepo) IsEnabled(ctx context.Context, conditions FeatureEnableQuery) bool {
	if conditions.FeatureID == uuid.Nil && conditions.FeatureName == "" {
		return false
	}
	//	Get Feature
	var f models.Feature
	query := r.db.DB().WithContext(ctx)
	if conditions.FeatureID != uuid.Nil {
		query = query.Where("id = @feature_id", sql.Named("feature_id", conditions.FeatureID))
	}
	if conditions.FeatureName != "" {
		query = query.Where("name = @name", sql.Named("name", conditions.FeatureName))
	}
	err := query.Find(&f).Error
	if err != nil {
		internal.Logger.ErrorContext(ctx, "failed to find feature", slog.Any("error", err))
		return false
	}
	if !f.IsEnabled() {
		internal.Logger.InfoContext(ctx, "feature is disabled")
		return false
	}
	var count int64

	// Get Feature's Strategies
	// FIXME: this cannot run correctly, need to query data in JSON
	if conditions.Key != "" {
		query = r.db.DB().Model(&models.Strategy{}).WithContext(ctx).
			Where("feature_id = ?", f.ID).
			Where(datatypes.JSONQuery("data").HasKey(conditions.Key)).
			Where(datatypes.JSONQuery("data").HasKey(conditions.Key, conditions.Value))

		if conditions.Environment != "" {
			query = query.Where(datatypes.JSONArrayQuery("environment").Contains(conditions.Environment))
		}
		query.Count(&count)
		if count == 0 {
			return false
		}
	}

	if conditions.CollectionValue != "" {
		query = r.db.DB().Model(&models.CollectionDetail{}).WithContext(ctx).
			Joins("inner join feature_collection on feature_collection.collection_id = collection_detail.collection_id").
			Where("feature_collection.feature_id = ?", f.ID).
			Where("collection_detail.value = ?", conditions.CollectionValue)
		query.Count(&count)
		if count == 0 {
			return false
		}
	}

	return true
}
