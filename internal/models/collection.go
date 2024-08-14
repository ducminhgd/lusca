package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Collection struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
}

func (Collection) TableName() string {
	return "collection"
}

func (m *Collection) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

type CollectionDetail struct {
	CollectionID uuid.UUID `json:"collection_id" gorm:"column:collection_id"`
	Value        string    `json:"value" gorm:"column:value"`
}

func (CollectionDetail) TableName() string {
	return "collection_detail"
}

type FeatureCollection struct {
	CollectionID uuid.UUID `json:"collection_id" gorm:"primaryKey;column:collection_id"`
	FeatureID    uuid.UUID `json:"feature_id" gorm:"primaryKey;column:feature_id"`
}

func (FeatureCollection) TableName() string {
	return "feature_collection"
}
