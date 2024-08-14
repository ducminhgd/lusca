package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	STRATEGYKV_TYPE_UNKNOWN = 0
	STRATEGYKV_TYPE_VALUE   = 1
	STRATEGYKV_TYPE_LIST    = 2
)

type StrategyKV struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey;column:id"`
	FeatureID   uuid.UUID      `json:"feature_id" gorm:"column:feature_id"`
	Environment datatypes.JSON `json:"environment" gorm:"column:environment;type:jsonb"`
	Type        int            `json:"type" gorm:"column:type"`
	Key         string         `json:"key" gorm:"column:key"`
	Value       datatypes.JSON `json:"value" gorm:"column:value;type:jsonb"`
}

func (StrategyKV) TableName() string {
	return "strategy_kv"
}

func (m *StrategyKV) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}
