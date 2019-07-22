package prometheus

import "testing"

func TestDefaultCounters(t *testing.T) {
	t.Run("First call", func(t *testing.T) {
		if counters.registered {
			t.Fatal("Counter must be not registered")
		}

		if counter := DefaultCounters(); !counter.registered {
			t.Fatal("Counter must be registered")
		}
	})

	t.Run("Registered flag check", func(t *testing.T) {
		// On duplicate prometheus throws the panic
		// This situation should be protected by registered flag
		defer func() {
			if rvr := recover(); rvr != nil {
				t.Fatal(rvr)
			}
		}()

		if counter := DefaultCounters(); !counter.registered {
			t.Fatal("Counter must be registered")
		}
	})
}
