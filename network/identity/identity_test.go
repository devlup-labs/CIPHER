package identity

import (
	"testing"
)

func TestLoadOrCreate(t *testing.T) {
	// Temporarily override user config dir for tests if possible.
	// Since os.UserConfigDir is a standard library function, we can't easily mock it without refactoring.
	// We'll test the actual behavior but in a controlled way if needed, or just let it run in the real environment.
	
	priv1, err := LoadOrCreate()
	if err != nil {
		t.Fatalf("Failed to create first key: %v", err)
	}

	priv2, err := LoadOrCreate()
	if err != nil {
		t.Fatalf("Failed to load second key: %v", err)
	}

	if !priv1.Equals(priv2) {
		t.Fatalf("Expected keys to be equal on reload")
	}
}
