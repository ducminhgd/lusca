package require

import (
	"log"
	"lusca/config"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type DBManager interface {
	NewDialector(dbAttrs config.DBAttributes) gorm.Dialector
	NewConnection(dbAttrs config.DBAttributes, gormConfig *gorm.Config) (*gorm.DB, error)
}

type dbManager struct {
	sources  []gorm.Dialector
	replicas []gorm.Dialector
	conn     *gorm.DB
}

func NewDatabaseManager() *dbManager {
	return &dbManager{
		sources:  []gorm.Dialector{},
		replicas: []gorm.Dialector{},
		conn:     &gorm.DB{},
	}
}

func NewDialector(dbAttrs config.DBAttributes) gorm.Dialector {
	return postgres.Open(dbAttrs.GetPostgresDSN())
}

func (m *dbManager) WithSourceDialector(dbAttrs config.DBAttributes) *dbManager {
	m.sources = append(m.sources, NewDialector(dbAttrs))
	return m
}

func (m *dbManager) WithReplicaDialector(dbAttrs config.DBAttributes) *dbManager {
	m.replicas = append(m.replicas, NewDialector(dbAttrs))
	return m
}

func (m *dbManager) Connect(dbConfig config.DBConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}
	gormConfig.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                        // Slow SQL threshold
			LogLevel:                  logger.LogLevel(dbConfig.LogLevel), // Log level
			IgnoreRecordNotFoundError: true,                               // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                               // Don't include params in the SQL log
			Colorful:                  true,                               // Disable color
		},
	)
	db, err := gorm.Open(m.sources[0], gormConfig)
	if err != nil {
		return nil, err
	}
	err = db.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:           m.sources,
			Replicas:          m.replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime).
		SetConnMaxLifetime(dbConfig.ConnMaxLifeTime).
		SetMaxIdleConns(dbConfig.MaxIdleConns).
		SetMaxOpenConns(dbConfig.MaxOpenConns),
	)

	if err != nil {
		return nil, err
	}
	return db, nil
}
