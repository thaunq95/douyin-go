package douyin

// BaseRequestParams represents standard query params needed for Douyin requests
type BaseRequestParams struct {
	DevicePlatform           string `url:"device_platform"`
	Aid                      string `url:"aid"`
	Channel                  string `url:"channel"`
	PcClientType             string `url:"pc_client_type"`
	VersionCode              string `url:"version_code"`
	VersionName              string `url:"version_name"`
	CookieEnabled            string `url:"cookie_enabled"`
	ScreenWidth              string `url:"screen_width"`
	ScreenHeight             string `url:"screen_height"`
	BrowserLanguage          string `url:"browser_language"`
	BrowserPlatform          string `url:"browser_platform"`
	BrowserName              string `url:"browser_name"`
	BrowserVersion           string `url:"browser_version"`
	BrowserOnline            string `url:"browser_online"`
	EngineName               string `url:"engine_name"`
	EngineVersion            string `url:"engine_version"`
	OsName                   string `url:"os_name"`
	OsVersion                string `url:"os_version"`
	CpuCoreNum               string `url:"cpu_core_num"`
	DeviceMemory             string `url:"device_memory"`
	Platform                 string `url:"platform"`
	Downlink                 string `url:"downlink"`
	EffectiveType            string `url:"effective_type"`
	FromUserPage             string `url:"from_user_page"`
	LocateQuery              string `url:"locate_query"`
	NeedTimeList             string `url:"need_time_list"`
	PcLibraDivert            string `url:"pc_libra_divert"`
	PublishVideoStrategyType string `url:"publish_video_strategy_type"`
	RoundTripTime            string `url:"round_trip_time"`
	ShowLiveReplayStrategy   string `url:"show_live_replay_strategy"`
	TimeListQuery            string `url:"time_list_query"`
	WhaleCutToken            string `url:"whale_cut_token"`
	UpdateVersionCode        string `url:"update_version_code"`
}

// DefaultParams returns a map of default query parameters for Douyin API requests
func DefaultParams() map[string]string {
	return map[string]string{
		"device_platform":             "webapp",
		"aid":                         "6383",
		"channel":                     "channel_pc_web",
		"pc_client_type":              "1",
		"version_code":                "290100",
		"version_name":                "29.1.0",
		"cookie_enabled":              "true",
		"screen_width":                "1920",
		"screen_height":               "1080",
		"browser_language":            "zh-CN",
		"browser_platform":            "Win32",
		"browser_name":                "Chrome",
		"browser_version":             "139.0.0.0",
		"browser_online":              "true",
		"engine_name":                 "Blink",
		"engine_version":              "139.0.0.0",
		"os_name":                     "Windows",
		"os_version":                  "10",
		"cpu_core_num":                "12",
		"device_memory":               "8",
		"platform":                    "PC",
		"downlink":                    "10",
		"effective_type":              "4g",
		"from_user_page":              "1",
		"locate_query":                "false",
		"need_time_list":              "1",
		"pc_libra_divert":             "Windows",
		"publish_video_strategy_type": "2",
		"round_trip_time":             "0",
		"show_live_replay_strategy":   "1",
		"time_list_query":             "0",
		"whale_cut_token":             "",
		"update_version_code":         "170400",
	}
}

// AuthorInfo represents user details within a post response
type AuthorInfo struct {
	UID       string `json:"uid"`
	SecUID    string `json:"sec_uid"`
	Nickname  string `json:"nickname"`
	Signature string `json:"signature"`
	ShortID   string `json:"short_id"`
	UniqueID  string `json:"unique_id"`
}

// ImageItem represents an image URL item in a slideshow post
type ImageItem struct {
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

// PlayAddr represents the playback address properties
type PlayAddr struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

// VideoInfo represents metadata about the video
type VideoInfo struct {
	PlayAddr    PlayAddr `json:"play_addr"`
	Cover       PlayAddr `json:"cover"`
	OriginCover PlayAddr `json:"origin_cover"`
	Duration    int64    `json:"duration"`
	Ratio       string   `json:"ratio"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
}

// MusicInfo represents audio details for a post
type MusicInfo struct {
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	PlayURL PlayAddr `json:"play_url"`
}

// StatisticsInfo represents likes, shares, comments counts
type StatisticsInfo struct {
	CommentCount int64 `json:"comment_count"`
	DiggCount    int64 `json:"digg_count"`
	ShareCount   int64 `json:"share_count"`
	PlayCount    int64 `json:"play_count"`
	CollectCount int64 `json:"collect_count"`
}

// AwemeDetail represents a single Douyin post (video or slideshow)
type AwemeDetail struct {
	AwemeID    string         `json:"aweme_id"`
	Desc       string         `json:"desc"`
	CreateTime int64          `json:"create_time"`
	Author     AuthorInfo     `json:"author"`
	Video      VideoInfo      `json:"video"`
	Images     []ImageItem    `json:"images"`
	Music      MusicInfo      `json:"music"`
	Statistics StatisticsInfo `json:"statistics"`
	AwemeType  int            `json:"aweme_type"` // 0: Video, 68: Image Slideshow
}

// PostDetailResponse represents the Douyin POST_DETAIL response schema
type PostDetailResponse struct {
	StatusCode int         `json:"status_code"`
	AwemeInfo  AwemeDetail `json:"aweme_detail"`
}

// UserPostResponse represents the Douyin USER_POST or USER_FAVORITE response schema
type UserPostResponse struct {
	StatusCode int           `json:"status_code"`
	AwemeList  []AwemeDetail `json:"aweme_list"`
	MaxCursor  int64         `json:"max_cursor"`
	HasMore    int           `json:"has_more"` // 1: has more, 0: no more
}

// UserProfileResponse represents the response containing user profile details
type UserProfileResponse struct {
	StatusCode int `json:"status_code"`
	User       struct {
		Nickname       string `json:"nickname"`
		SecUID         string `json:"sec_uid"`
		UniqueID       string `json:"unique_id"`
		ShortID        string `json:"short_id"`
		Signature      string `json:"signature"`
		FollowingCount int64  `json:"following_count"`
		FollowerCount  int64  `json:"follower_count"`
		TotalFavorited int64  `json:"total_favorited"`
		AwemeCount     int64  `json:"aweme_count"`
	} `json:"user"`
}
