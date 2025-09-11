package tests

import (
	"radius-server/src/common/logger"
	"testing"
)

func RunTestNegativeCases(t *testing.T) {
	logger.Logger.Info().Msg("Running TestNegativeCases")
	t.Run("should pass", func(t *testing.T) {
		if 1+1 != 2 {
			t.Fatalf("math is broken")
		}
	})
}
