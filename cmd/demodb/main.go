package main

import (
	"context"
	"fmt"

	"github.com/ducminhgd/lusca/config"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/ducminhgd/lusca/internal/repositories"
	"github.com/ducminhgd/lusca/internal/require"
	"gorm.io/datatypes"
	"gorm.io/plugin/dbresolver"
)

func main() {
	var (
		cfg = config.Load()
	)

	// createSampleData(cfg)
	sampleQuery(cfg)
}

func sampleQuery(cfg config.Config) {
	dbManager, err := require.NewDatabaseManager(cfg.Databases)
	if err != nil {
		panic(err)
	}
	featureRepo := repositories.NewFeatureRepo(dbManager)
	result := featureRepo.IsEnabled(context.Background(), repositories.FeatureEnableQuery{
		FeatureName: "feature-1",
		Key:         "region",
		Value:       "vn",
	})
	fmt.Println(result)
}

func createSampleData(cfg config.Config) {
	dbManager, err := require.NewDatabaseManager(cfg.Databases)
	if err != nil {
		panic(err)
	}
	db := dbManager.DB()
	f := models.Feature{
		Name:        "feature-1",
		Description: "feature-1 description",
		Status:      2,
	}
	db.Clauses(dbresolver.Write).Create(&f)
	fmt.Println(f)

	s := models.Strategy{
		FeatureID:   f.ID,
		Environment: datatypes.JSON([]byte(`["dev"]`)),
		Data:        datatypes.JSON([]byte(`{"region":"vn"}`)),
	}
	db.Clauses(dbresolver.Write).Create(&s)
	fmt.Println(s)
}
