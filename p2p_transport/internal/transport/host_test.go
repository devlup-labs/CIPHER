package transport

import (
	"context"
	"testing"
)

func TestNewNode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host, err := NewNode(ctx, 0, nil, "", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if host == nil {
		t.Fatalf("Expected a host, got nil")
	}
	
	if len(host.Addrs()) == 0 {
		t.Fatalf("Expected at least one listen address")
	}

	host.Close()
}
