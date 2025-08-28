package database

import (
	"database/sql"
	"radius-server/src/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "gorm.io/driver/postgres"
)

func RunMigrations() error {
	dsn := config.AppConfig.Database.Dsn
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	migration, err := migrate.NewWithDatabaseInstance(
		"file://./src/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
