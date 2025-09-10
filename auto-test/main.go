package main

import (
	"fmt"
	"radius-server/auto-test/services"
	"radius-server/src/common/logger"
	"radius-server/src/config"
	"radius-server/src/database"
	"time"
)

func init() {
	logger.InitializeLogger()
}

func main() {
	config.LoadConfig()

	if err := database.Connect(); err != nil {
		logger.Logger.Fatal().Msgf("Connection to database error. %s", err.Error())
	}

	err := services.CreateNas()
	if err != nil {
		logger.Logger.Fatal().Err(err)
	}

	fmt.Println(services.Nas)

	time.Sleep(time.Second * 10)
	err = services.DeleteNas()
	if err != nil {
		logger.Logger.Fatal().Err(err)
	}
}
