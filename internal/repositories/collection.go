package repositories

import (
	"context"
	"database/sql"

	"github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectionDetailQuery struct {
	CollectionID uuid.UUID `json:"collection_id" default:""`
	Value        string    `json:"value" default:""`
	Value_Like   string    `json:"value__like" default:""`
	Value_In     []string  `json:"value__in"`
}

type collectionDetailRepo struct {
	db *db.GORMManager
}

func NewCollectionDetailRepo(db *db.GORMManager) *collectionDetailRepo {
	return &collectionDetailRepo{
		db: db,
	}
}

func (r *collectionDetailRepo) Upsert(ctx context.Context, m *models.CollectionDetail) error {
	return r.db.DB().WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error
}

func (r *collectionDetailRepo) Detele(ctx context.Context, conditions CollectionDetailQuery) (int, error) {
	q := r.db.DB().WithContext(ctx).Model(&models.CollectionDetail{})
	if conditions.CollectionID != uuid.Nil {
		q = q.Where("collection_id = @collectionID", sql.Named("collectionID", conditions.CollectionID))
	}
	if conditions.Value != "" {
		q = q.Where("value = @value", sql.Named("value", conditions.Value))
	}
	if conditions.Value_Like != "" {
		q = q.Where("value like @value__like", sql.Named("value__like", "%"+conditions.Value_Like+"%"))
	}
	if len(conditions.Value_In) > 0 {
		q = q.Where("value in (@value__in)", sql.Named("value__in", conditions.Value_In))
	}
	return int(q.Delete(&models.CollectionDetail{}).RowsAffected), q.Error
}

func (r *collectionDetailRepo) GetByCollectionID(ctx context.Context, id uuid.UUID, opts QueryOptions) ([]models.CollectionDetail, error) {
	var cds []models.CollectionDetail
	q := r.db.DB().WithContext(ctx).Model(&models.CollectionDetail{}).Where("collection_id = ?", id)
	err := q.Limit(opts.Limit).Offset(opts.Offset).Find(&cds).Error
	return cds, err
}

type CollectionRepo struct {
	db *db.GORMManager
}

type CollectionQuery struct {
	Name_Like string `json:"name__like" default:""`
}

func NewCollectionRepo(db *db.GORMManager) *CollectionRepo {
	return &CollectionRepo{
		db: db,
	}
}

func (r *CollectionRepo) Create(ctx context.Context, m *models.Collection) (int, error) {
	cmd := r.db.DB().WithContext(ctx)
	if len(m.Details) == 0 {
		cmd = cmd.Omit("Details")
	}
	result := cmd.Create(&m)
	return int(result.RowsAffected), result.Error
}

func (r *CollectionRepo) Update(ctx context.Context, m *models.Collection) (int, error) {
	result := r.db.DB().WithContext(ctx).Model(&models.Collection{}).Where("id = ?", m.ID).Updates(&m)
	return int(result.RowsAffected), result.Error
}

func (r *CollectionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	var c models.Collection
	err := r.db.DB().Preload("Details").WithContext(ctx).Model(&models.Collection{}).Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *CollectionRepo) GetByName(ctx context.Context, name string) (*models.Collection, error) {
	var c models.Collection
	err := r.db.DB().Preload("Details").WithContext(ctx).Model(&models.Collection{}).Where("name = ?", name).First(&c).Error
	return &c, err
}

func (r *CollectionRepo) GetList(ctx context.Context, conditions CollectionQuery, opts QueryOptions) ([]models.Collection, int64, error) {
	var (
		cs    []models.Collection
		count int64 = 0
	)
	q := r.db.DB().Preload("Details").WithContext(ctx).Model(&models.Collection{})
	if conditions.Name_Like != "" {
		q = q.Where("name like ?", "%"+conditions.Name_Like+"%")
	}
	q = q.Session(&gorm.Session{})
	result := q.Count(&count)
	if result.Error != nil {
		return cs, count, result.Error
	}
	result = q.Limit(opts.Limit).Offset(opts.Offset).Find(&cs)
	if result.Error != nil {
		return cs, count, result.Error
	}
	return cs, count, result.Error
}

func (r *CollectionRepo) LinkFeatureCollection(ctx context.Context, featureID, collectionID uuid.UUID) error {
	return r.db.DB().WithContext(ctx).Create(&models.FeatureCollection{
		FeatureID:    featureID,
		CollectionID: collectionID,
	}).Error
}
