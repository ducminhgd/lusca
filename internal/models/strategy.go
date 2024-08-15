package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Strategy struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey;column:id"`
	FeatureID   uuid.UUID      `json:"feature_id" gorm:"column:feature_id"`
	Environment datatypes.JSON `json:"environment" gorm:"column:environment;type:jsonb"`
	Data        datatypes.JSON `json:"data" gorm:"column:data;type:jsonb"`
}

func (Strategy) TableName() string {
	return "strategy"
}

func (m *Strategy) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}
