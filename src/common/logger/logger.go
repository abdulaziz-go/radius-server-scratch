package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitializeLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.FormatLevel = func(i interface{}) string {
		switch i {
		case "info":
			return fmt.Sprintf("\033[32m%s\033[0m", i) // Green for info
		case "warn":
			return fmt.Sprintf("\033[33m%s\033[0m", i) // Yellow for warn
		case "error":
			return fmt.Sprintf("\033[31m%s\033[0m", i) // Red for error
		case "fatal":
			return fmt.Sprintf("\033[41m%s\033[0m", i) // White text on red background for fatal
		default:
			return fmt.Sprintf("\033[37m%s\033[0m", i) // Default color for other levels
		}
	}
	Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}

func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		start := time.Now()
		requestID := c.Get("X-Request-ID")
		event := Logger.Info().
			Str("requestid", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", time.Since(start))
		if err != nil {
			event.Msg("Request failed. Error - " + err.Error())
		} else {
			event.Msg("Request completed")
		}
		return err
	}
}

func SetDebugLevel() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
