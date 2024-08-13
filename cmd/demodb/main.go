package main

import (
	"fmt"
	"lusca/config"
	"lusca/internal/models"
	"lusca/internal/require"

	"github.com/google/uuid"
	"gorm.io/plugin/dbresolver"
)

func main() {
	var (
		cfg = config.Load()
	)
	dbManager := require.NewDatabaseManager()
	dbManager.WithSourceDialector(cfg.Databases.Source)
	if cfg.Databases.UseReplication {
		dbManager.WithReplicaDialector(cfg.Databases.Replicas[0])
	}
	db, _ := dbManager.Connect(cfg.Databases, nil)

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
