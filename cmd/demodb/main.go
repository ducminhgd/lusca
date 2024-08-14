package main

import (
	"fmt"

	"github.com/ducminhgd/lusca/config"
	"github.com/ducminhgd/lusca/internal/models"
	"github.com/ducminhgd/lusca/internal/require"
	"github.com/google/uuid"
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

	uuidv7, _ := uuid.NewV7()
	fmt.Println(uuidv7)
	f := models.Feature{
		ID:          uuidv7,
		Name:        "feature-2",
		Description: "feature-2 description",
		Status:      1,
	}
	db.Clauses(dbresolver.Write).Create(&f)
}
