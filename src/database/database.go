package database

import (
	"radius-server/src/config"
	timeUtil "radius-server/src/utils/time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// Database instance
var DbConn *gorm.DB

func Connect() error {
	var err error
	// Use DSN string to open
	DbConn, err = gorm.Open(postgres.Open(config.AppConfig.Database.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DbConn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(config.AppConfig.Database.Dsn)},
		Replicas: []gorm.Dialector{},
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		return err
	}
	sqlDB, err := DbConn.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.Connection.MaxNumber)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.Connection.OpenMaxNumber)
	sqlDB.SetConnMaxLifetime(timeUtil.DurationSeconds(config.AppConfig.Database.Connection.MaxLifetimeSec))
	sqlDB.SetConnMaxIdleTime(timeUtil.DurationSeconds(config.AppConfig.Database.Connection.MaxIdleTimeSec))
	return nil
}

func getDbLogger() logger.Interface {
	var dbLogger logger.Interface
	if config.AppConfig.IsDebug && config.AppConfig.Database.Logger {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}
	return dbLogger
}
