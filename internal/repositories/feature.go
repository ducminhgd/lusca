package repository

import (
	"github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/google/uuid"
)

type featureRepo struct {
	db *db.GORMManager
}

func NewFeatureRepo(db *db.GORMManager) *featureRepo {
	return &featureRepo{
		db: db,
	}
}

func (r *featureRepo) GetByID(id uuid.UUID) (*models.Feature, error) {
	var f models.Feature
	err := r.db.DB().Model(&models.Feature{}).Where("id = ?", id).First(&f).Error
	return &f, err
}

func (r *featureRepo) Create(name string, description string, status int) (*models.Feature, error) {
	f := models.Feature{
		Name:        name,
		Description: description,
		Status:      status,
	}
	err := r.db.DB().Create(&f).Error
	return &f, err
}

func (r *featureRepo) Update(id uuid.UUID, name string, description string, status int) (*models.Feature, error) {
	f := models.Feature{
		Name:        name,
		Description: description,
		Status:      status,
	}
	err := r.db.DB().Model(&models.Feature{}).Where("id = ?", id).Updates(&f).Error
	return &f, err
}

func (r *featureRepo) GetList(name string, description string, status int, limit int, offset int) ([]models.Feature, error) {
	var fs []models.Feature
	err := r.db.DB().Model(&models.Feature{}).
		Where("name like ?", "%"+name+"%").
		Where("description like ?", "%"+description+"%").
		Where("status = ?", status).
		Order("created_at desc").
		Limit(limit).Offset(offset).Find(&fs).Error
	return fs, err
}
