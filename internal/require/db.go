package require

import (
	gao_db "github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager interface {
	NewPostgresDialector(config.DBAttributes) (gorm.Dialector, error)
	NewDatabaseManager(config.DBConfig) (*gao_db.GORMManager, error)
}

func NewPostgresDialector(cfg config.DBAttributes) (gorm.Dialector, error) {
	dsn := cfg.GetPostgresDSN()
	return postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), nil
}

func NewDatabaseManager(cfg config.DBConfig) (*gao_db.GORMManager, error) {
	prm, err := NewPostgresDialector(cfg.Source)
	if err != nil {
		return nil, err
	}

	m, _ := gao_db.NewGORMManager(prm, &gorm.Config{})

	if cfg.UseReplication {
		rpl, err := NewPostgresDialector(cfg.Replica)
		if err != nil {
			return nil, err
		}

		m, _ = m.AddReplicaDialector(rpl, gao_db.PoolConfig{})
	}
	return m, nil
}
