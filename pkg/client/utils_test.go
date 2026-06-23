package client

import (
	"testing"
)

func TestExtractURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"Check out this video: https://www.douyin.com/video/7345492945006595379 amazing!", "https://www.douyin.com/video/7345492945006595379", false},
		{"https://v.douyin.com/abcde/", "https://v.douyin.com/abcde/", false},
		{"no url here", "", true},
	}

	for _, tc := range tests {
		got, err := ExtractURL(tc.input)
		if (err != nil) != tc.err {
			t.Errorf("ExtractURL(%q) error state got %v, expected %v", tc.input, err != nil, tc.err)
		}
		if got != tc.expected {
			t.Errorf("ExtractURL(%q) got %q, expected %q", tc.input, got, tc.expected)
		}
	}
}

func TestExtractSecUID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"https://www.douyin.com/user/MS4wLjABAAAAEs86TBQPNwAo-RGrcxWyCdwKhI66AK3Pqf3ieo6HaxI?showTab=post", "MS4wLjABAAAAEs86TBQPNwAo-RGrcxWyCdwKhI66AK3Pqf3ieo6HaxI", false},
		{"https://www.douyin.com/discover?sec_uid=MS4wLjABAAAA12345", "MS4wLjABAAAA12345", false},
	}

	for _, tc := range tests {
		got, err := ExtractSecUID(tc.input)
		if (err != nil) != tc.err {
			t.Errorf("ExtractSecUID(%q) error state got %v, expected %v", tc.input, err != nil, tc.err)
		}
		if got != tc.expected {
			t.Errorf("ExtractSecUID(%q) got %q, expected %q", tc.input, got, tc.expected)
		}
	}
}

func TestExtractAwemeID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"https://www.douyin.com/video/7345492945006595379", "7345492945006595379", false},
		{"https://www.douyin.com/note/12345", "12345", false},
		{"https://www.douyin.com/discover?modal_id=67890", "67890", false},
		{"https://www.douyin.com/web/api/v2/aweme/like/?vid=98765", "98765", false},
		{"https://www.douyin.com/123456789", "123456789", false}, // path only fallback
	}

	for _, tc := range tests {
		got, err := ExtractAwemeID(tc.input)
		if (err != nil) != tc.err {
			t.Errorf("ExtractAwemeID(%q) error state got %v, expected %v", tc.input, err != nil, tc.err)
		}
		if got != tc.expected {
			t.Errorf("ExtractAwemeID(%q) got %q, expected %q", tc.input, got, tc.expected)
		}
	}
}
