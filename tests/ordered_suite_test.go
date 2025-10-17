package tests

import "testing"

func TestOrderedSuite(t *testing.T) {
	t.Run("Auth", func(t *testing.T) {
		RunTestAuth(t)
	})
	t.Run("Accounting", func(t *testing.T) {
		RunTestAccounting(t)
	})
	t.Run("CoaDisconnect", func(t *testing.T) {
		RunTestCoaDisconnect(t)
	})
	t.Run("NegativeCases", func(t *testing.T) {
		RunTestNegativeCases(t)
	})
	t.Run("Policy", func(t *testing.T) {
		RunTestPolicy(t)
	})
}
