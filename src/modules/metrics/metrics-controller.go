package metricsModule

import (
	"radius-server/src/metrics"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetMetrics(c *fiber.Ctx) error {
	metrics := metrics.GetMetricsPromtheusFormatted()
	c.Set("Content-Type", "text/plain; version=0.0.4")

	var builder strings.Builder
	for i := range metrics {
		builder.WriteString(metrics[i])
	}

	return c.SendString(builder.String())
}
