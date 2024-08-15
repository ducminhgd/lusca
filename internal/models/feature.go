package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	FEATURE_STATUS_UNKNOWN  = 0
	FEATURE_STATUS_DISABLED = 1
	FEATURE_STATUS_ENABLED  = 2
)

type Feature struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Status      int       `json:"status" gorm:"column:status"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (Feature) TableName() string {
	return "feature"
}

func (m *Feature) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

func (m *Feature) IsEnabled() bool {
	return m.Status == FEATURE_STATUS_ENABLED
}
