package tests

import (
	"radius-server/src/common/logger"
	"testing"
)

// Example test template. Duplicate and adapt for new tests.
func RunTestExample(t *testing.T) {
	logger.Logger.Info().Msg("Running TestExample")
	t.Run("should pass", func(t *testing.T) {
		if 1+1 != 2 {
			t.Fatalf("math is broken")
		}
	})
}
