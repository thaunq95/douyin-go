package douyin

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thaunq95/douyin-go/pkg/client"
)

type Crawler struct {
	Client *client.Client
}

// NewCrawler creates a new crawler instance. It fetches a ttwid if no cookie is provided.
func NewCrawler(cookie string, proxy string) (*Crawler, error) {
	c, err := client.NewClient(cookie, proxy)
	if err != nil {
		return nil, err
	}

	// If no cookie is set, fetch a ttwid to populate client's request credentials
	if cookie == "" {
		ttwid, err := c.FetchTTWid()
		if err == nil && ttwid != "" {
			c.Cookie = fmt.Sprintf("ttwid=%s;", ttwid)
		}
	}

	return &Crawler{
		Client: c,
	}, nil
}

// GetVideoDetail fetches details of a single post by its aweme_id
func (cr *Crawler) GetVideoDetail(awemeID string) (*AwemeDetail, error) {
	params := DefaultParams()
	params["aweme_id"] = awemeID

	data, err := cr.Client.Get(PostDetailURL, params)
	if err != nil {
		return nil, err
	}

	var resp PostDetailResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if resp.StatusCode != 0 {
		return nil, fmt.Errorf("douyin api returned status code %d", resp.StatusCode)
	}

	if resp.AwemeInfo.AwemeID == "" {
		return nil, fmt.Errorf("empty aweme info returned")
	}

	return &resp.AwemeInfo, nil
}

// GetUserPosts fetches post lists published by a user
func (cr *Crawler) GetUserPosts(secUserID string, maxCursor int64, count int) (*UserPostResponse, error) {
	params := DefaultParams()
	params["sec_user_id"] = secUserID
	params["max_cursor"] = strconv.FormatInt(maxCursor, 10)
	params["count"] = strconv.Itoa(count)

	data, err := cr.Client.Get(UserPostURL, params)
	if err != nil {
		return nil, err
	}

	var resp UserPostResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &resp, nil
}

// GetUserLikes fetches post lists liked by a user
func (cr *Crawler) GetUserLikes(secUserID string, maxCursor int64, count int) (*UserPostResponse, error) {
	params := DefaultParams()
	params["sec_user_id"] = secUserID
	params["max_cursor"] = strconv.FormatInt(maxCursor, 10)
	params["count"] = strconv.Itoa(count)

	data, err := cr.Client.Get(UserLikeURL, params)
	if err != nil {
		return nil, err
	}

	var resp UserPostResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &resp, nil
}

// GetUserProfile fetches profile details of a user
func (cr *Crawler) GetUserProfile(secUserID string) (*UserProfileResponse, error) {
	params := DefaultParams()
	params["sec_user_id"] = secUserID

	data, err := cr.Client.Get(UserDetailURL, params)
	if err != nil {
		return nil, err
	}

	var resp UserProfileResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &resp, nil
}
