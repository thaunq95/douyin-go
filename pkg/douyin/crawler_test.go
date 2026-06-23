package douyin

import (
	"testing"
)

func TestNewCrawler(t *testing.T) {
	cr, err := NewCrawler("", "")
	if err != nil {
		t.Fatalf("failed to create crawler: %v", err)
	}

	if cr.Client == nil {
		t.Errorf("Client should not be nil")
	}
}

func TestGetVideoDetail(t *testing.T) {
	cr, err := NewCrawler("", "")
	if err != nil {
		t.Fatalf("failed to create crawler: %v", err)
	}

	// Known public video ID
	awemeID := "7345492945006595379"
	detail, err := cr.GetVideoDetail(awemeID)
	if err != nil {
		t.Skipf("Skipping integration test (requires internet / valid cookie or IP blocks): %v", err)
		return
	}

	if detail.AwemeID != awemeID {
		t.Errorf("expected aweme_id %s, got %s", awemeID, detail.AwemeID)
	}

	t.Logf("Success! Found video title: %s", detail.Desc)
}
