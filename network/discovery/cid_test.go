package discovery

import (
	"testing"

	"cipher/network/content/core"
)

func TestContentIDToCID(t *testing.T) {
	var id core.ContentID

	for i := range id {
		id[i] = byte(i)
	}

	c, err := contentIDToCID(id)
	if err != nil {
		t.Fatalf("contentIDToCID failed: %v", err)
	}

	if !c.Defined() {
		t.Fatal("expected defined CID")
	}

	if c.Version() != 1 {
		t.Fatalf("expected CIDv1, got CIDv%d", c.Version())
	}
}
