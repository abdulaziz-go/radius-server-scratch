package tests

import (
	"radius-server/src/common/logger"
	"testing"
)

func RunTestAccounting(t *testing.T) {
	logger.Logger.Info().Msg("Running TestAccounting")
	t.Run("should pass", func(t *testing.T) {
		if 1+1 != 2 {
			t.Fatalf("math is broken")
		}
	})
}
