package repositories

import (
	"context"

	"github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/google/uuid"
)

type StrategyQuery struct {
	ID_In          []uuid.UUID `json:"id__in"`
	FeatureID      uuid.UUID   `json:"feature_id" default:""`
	Environment    string      `json:"environment" default:""`
	Environment_In []string    `json:"environment__in"`
}

type strategyRepo struct {
	db *db.GORMManager
}

func NewStrategyKVRepo(db *db.GORMManager) *strategyRepo {
	return &strategyRepo{
		db: db,
	}
}

func (r *strategyRepo) GetList(ctx context.Context, conditions StrategyQuery, opts QueryOptions) ([]models.Strategy, error) {
	var ss []models.Strategy
	q := r.db.DB().WithContext(ctx).Model(&models.Strategy{})
	if len(conditions.ID_In) > 0 {
		q = q.Where("id in (?)", conditions.ID_In)
	}
	if conditions.FeatureID != uuid.Nil {
		q = q.Where("feature_id = ?", conditions.FeatureID)
	}
	if conditions.Environment != "" {
		q = q.Where("environment = ?", conditions.Environment)
	}
	if len(conditions.Environment_In) > 0 {
		q = q.Where("environment in (?)", conditions.Environment_In)
	}
	err := q.Order("created_at desc").
		Limit(opts.Limit).Offset(opts.Offset).Find(&ss).Error
	return ss, err
}

func (r *strategyRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Strategy, error) {
	var s models.Strategy
	err := r.db.DB().WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", id).First(&s).Error
	return &s, err
}

func (r *strategyRepo) GetByFeatureID(ctx context.Context, featureID uuid.UUID) ([]models.Strategy, error) {
	var ss []models.Strategy
	err := r.db.DB().WithContext(ctx).Model(&models.Strategy{}).Where("feature_id = ?", featureID).Find(&ss).Error
	return ss, err
}

func (r *strategyRepo) Create(ctx context.Context, m *models.Strategy) (*models.Strategy, error) {
	err := r.db.DB().WithContext(ctx).Create(&m).Error
	return m, err
}
func (r *strategyRepo) Update(ctx context.Context, m *models.Strategy) (*models.Strategy, error) {
	err := r.db.DB().WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", m.ID).Updates(&m).Error
	return m, err
}
