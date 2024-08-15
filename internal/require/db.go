package require

import (
	"log"
	"os"
	"time"

	gao_db "github.com/ducminhgd/gao/db"
	"github.com/ducminhgd/lusca/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	m.WithLogger(logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,                   // Slow SQL threshold
			LogLevel:                  logger.LogLevel(cfg.LogLevel), // Log level
			IgnoreRecordNotFoundError: true,                          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                          // Don't include params in the SQL log
			Colorful:                  false,                         // Disable color
		},
	))

	if cfg.UseReplication {
		rpl, err := NewPostgresDialector(cfg.Replica)
		if err != nil {
			return nil, err
		}

		m, _ = m.AddReplicaDialector(rpl, gao_db.PoolConfig{})
	}
	return m, nil
}
