package main

import (
	"radius-server/src/common/logger"
	"radius-server/src/config"
	"radius-server/src/database"
	"radius-server/src/radius"
	"radius-server/src/redis"
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

	if err := redis.Connect(); err != nil {
		logger.Logger.Fatal().Msgf("Connection to redis error. %s", err.Error())
	}

	go func() {
		app, listenAddress := routes.New()
		if err := app.Listen(listenAddress); err != nil {
			logger.Logger.Fatal().Msgf("Startup error. %s", err.Error())
		}
	}()
	if err := radius.New().Start(); err != nil {
		logger.Logger.Fatal().Msgf("Startup radius server error. %s", err.Error())
	}
}
