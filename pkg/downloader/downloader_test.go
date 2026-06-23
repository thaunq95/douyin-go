package downloader

import (
	"testing"

	"github.com/thaunq95/douyin-go/pkg/douyin"
)

func TestFormatPath(t *testing.T) {
	detail := &douyin.AwemeDetail{
		AwemeID:    "123456",
		Desc:       "My Cool Video Description",
		CreateTime: 1718200000, // 2024-06-12 13:46:40 UTC
		Author: douyin.AuthorInfo{
			Nickname: "ThauNQ",
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{"{nickname}/{publish_date}_{title}", "ThauNQ/2024-06-12_My Cool Video Description"},
		{"{nickname}/likes/{publish_date}_{title}", "ThauNQ/likes/2024-06-12_My Cool Video Description"},
		{"{nickname}_posts/{title}_{aweme_id}", "ThauNQ_posts/My Cool Video Description_123456"},
		{"{nickname}/{title}/{publish_date}", "ThauNQ/My Cool Video Description/2024-06-12"},
	}

	for _, tc := range tests {
		got := FormatPath(tc.template, detail)
		if got != tc.expected {
			t.Errorf("FormatPath(%q) got %q, expected %q", tc.template, got, tc.expected)
		}
	}
}
