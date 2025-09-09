package main

import (
	"radius-server/src/common/logger"
	"radius-server/src/config"
	"radius-server/src/database"
	"radius-server/src/routes"
)

func init() {
	logger.InitializeLogger()
}

func main() {
	config.LoadConfig()

	if err := database.Connect(); err != nil {
		logger.Logger.Fatal().Msgf("Connection to database error. %s", err.Error())
	}
	if config.AppConfig.Database.AutoRunMigration {
		if err := database.RunMigrations(); err != nil {
			logger.Logger.Fatal().Msgf("Run migration error. %s", err.Error())
		}
	}

	app, listenAddress := routes.New()
	if err := app.Listen(listenAddress); err != nil {
		logger.Logger.Fatal().Msgf("Startup error. %s", err.Error())
	}
}
