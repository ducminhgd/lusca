package repositories

import (
	"context"

	"github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/google/uuid"
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
		q = q.Where("collection_id = ?", conditions.CollectionID)
	}
	if conditions.Value != "" {
		q = q.Where("value = ?", conditions.Value)
	}
	if conditions.Value_Like != "" {
		q = q.Where("value like ?", "%"+conditions.Value_Like+"%")
	}
	if len(conditions.Value_In) > 0 {
		q = q.Where("value in (?)", conditions.Value_In)
	}
	return int(q.Delete(&models.CollectionDetail{}).RowsAffected), q.Error
}

func (r *collectionDetailRepo) GetByCollectionID(ctx context.Context, id uuid.UUID, opts QueryOptions) ([]models.CollectionDetail, error) {
	var cds []models.CollectionDetail
	q := r.db.DB().WithContext(ctx).Model(&models.CollectionDetail{}).Where("collection_id = ?", id)
	err := q.Limit(opts.Limit).Offset(opts.Offset).Find(&cds).Error
	return cds, err
}

type collectionRepo struct {
	db *db.GORMManager
}

func NewCollectionRepo(db *db.GORMManager) *collectionRepo {
	return &collectionRepo{
		db: db,
	}
}

func (r *collectionRepo) Create(ctx context.Context, m *models.Collection) (*models.Collection, error) {
	err := r.db.DB().WithContext(ctx).Create(&m).Error
	return m, err
}

func (r *collectionRepo) Update(ctx context.Context, m *models.Collection) (*models.Collection, error) {
	err := r.db.DB().WithContext(ctx).Model(&models.Collection{}).Where("id = ?", m.ID).Updates(&m).Error
	return m, err
}

func (r *collectionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	var c models.Collection
	err := r.db.DB().WithContext(ctx).Model(&models.Collection{}).Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *collectionRepo) GetByName(ctx context.Context, name string) (*models.Collection, error) {
	var c models.Collection
	err := r.db.DB().WithContext(ctx).Model(&models.Collection{}).Where("name = ?", name).First(&c).Error
	return &c, err
}

func (r *collectionRepo) LinkFeatureCollection(ctx context.Context, featureID, collectionID uuid.UUID) error {
	return r.db.DB().WithContext(ctx).Create(&models.FeatureCollection{
		FeatureID:    featureID,
		CollectionID: collectionID,
	}).Error
}
