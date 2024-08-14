package main

import (
	"github.com/ducminhgd/lusca/config"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/ducminhgd/lusca/internal/require"
	"gorm.io/plugin/dbresolver"
)

func main() {
	var (
		cfg = config.Load()
	)
	dbManager, err := require.NewDatabaseManager(cfg.Databases)
	if err != nil {
		panic(err)
	}
	db := dbManager.DB()
	f := models.Feature{
		Name:        "feature-2",
		Description: "feature-2 description",
		Status:      1,
	}
	db.Clauses(dbresolver.Write).Create(&f)
}
