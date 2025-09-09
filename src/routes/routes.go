package routes

import (
	"radius-server/src/common/logger"
	"radius-server/src/config"
	apiModule "radius-server/src/modules/api"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func New() (*fiber.App, string) {
	// create app
	app := fiber.New(fiber.Config{
		AppName: config.AppConfig.AppName,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.Cors.Origins,
		AllowMethods:     config.AppConfig.Cors.Methods,
		AllowHeaders:     config.AppConfig.Cors.Headers,
		AllowCredentials: config.AppConfig.Cors.Credentials,
	}))
	app.Use(logger.LogMiddleware())

	apiMethods := app.Group("/")
	apiMethods.Get("/healthcheck", apiModule.HealthCheck)

	return app, ":" + strconv.Itoa(config.AppConfig.ServerPort)
}
