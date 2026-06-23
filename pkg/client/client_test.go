package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("", "")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if c.UserAgent == "" {
		t.Errorf("UserAgent should not be empty")
	}
}

func TestFetchTTWid(t *testing.T) {
	c, err := NewClient("", "")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// This is an integration test and depends on network access
	ttwid, err := c.FetchTTWid()
	if err != nil {
		t.Skipf("Skipping ttwid fetch test (requires internet / may fail due to blocks): %v", err)
		return
	}

	if ttwid == "" {
		t.Errorf("expected non-empty ttwid")
	} else {
		t.Logf("Successfully fetched ttwid: %s", ttwid)
	}
}
