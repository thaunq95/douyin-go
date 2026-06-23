package handlers

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/thaunq95/douyin-go/pkg/client"
	"github.com/thaunq95/douyin-go/pkg/config"
	"github.com/thaunq95/douyin-go/pkg/douyin"
	"github.com/thaunq95/douyin-go/pkg/downloader"
)

// Options holds CLI flag override parameters
type Options struct {
	Cookie      string
	Output      string
	Template    string
	Proxy       string
	Concurrency int
	Count       int
}

// Resolve merges configuration file settings with command line flag overrides
func Resolve(cfg *config.Config, opts Options) Options {
	resolved := opts
	if resolved.Cookie == "" {
		resolved.Cookie = cfg.Cookie
	}
	if resolved.Output == "" {
		resolved.Output = cfg.OutputDir
	}
	if resolved.Template == "" {
		resolved.Template = cfg.NamingTemplate
	}
	if resolved.Proxy == "" {
		resolved.Proxy = cfg.Proxy
	}
	if resolved.Concurrency <= 0 {
		resolved.Concurrency = cfg.Concurrency
	}
	if resolved.Concurrency <= 0 {
		resolved.Concurrency = 5 // default fallback
	}
	if resolved.Count <= 0 {
		resolved.Count = cfg.Count
	}
	if resolved.Count <= 0 {
		resolved.Count = 50 // default fallback
	}
	return resolved
}

// HandleVideo downloads a single video/post by URL
func HandleVideo(urlArg string, cfg *config.Config, opts Options) {
	resolved := Resolve(cfg, opts)

	fmt.Println("Initializing Douyin Downloader...")
	cr, err := douyin.NewCrawler(resolved.Cookie, resolved.Proxy)
	if err != nil {
		fmt.Printf("Error creating crawler: %v\n", err)
		os.Exit(1)
	}
	if cfg.UserAgent != "" {
		cr.Client.UserAgent = cfg.UserAgent
	}

	if err := os.MkdirAll(resolved.Output, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Parsing video/post URL...")
	awemeID, err := client.ExtractAwemeID(urlArg)
	if err != nil {
		fmt.Printf("Error extracting aweme_id: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found aweme_id: %s. Fetching details...\n", awemeID)

	detail, err := cr.GetVideoDetail(awemeID)
	if err != nil {
		fmt.Printf("Error fetching video details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found: %s (by %s)\n", detail.Desc, detail.Author.Nickname)
	if err := downloader.DownloadPost(detail, resolved.Output, resolved.Template, true); err != nil {
		fmt.Printf("Error downloading content: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Download completed successfully!")
}

// HandleUser downloads all posts published by a user profile
func HandleUser(userArg string, cfg *config.Config, opts Options) {
	resolved := Resolve(cfg, opts)

	fmt.Println("Initializing Douyin Downloader...")
	cr, err := douyin.NewCrawler(resolved.Cookie, resolved.Proxy)
	if err != nil {
		fmt.Printf("Error creating crawler: %v\n", err)
		os.Exit(1)
	}
	if cfg.UserAgent != "" {
		cr.Client.UserAgent = cfg.UserAgent
	}

	if err := os.MkdirAll(resolved.Output, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Parsing user profile URL...")
	secUID, err := client.ExtractSecUID(userArg)
	if err != nil {
		fmt.Printf("Error extracting sec_user_id: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found sec_user_id: %s. Fetching profile info...\n", secUID)

	profile, err := cr.GetUserProfile(secUID)
	if err != nil {
		fmt.Printf("Error fetching user profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("User profile: %s (ID: %s)\n", profile.User.Nickname, profile.User.UniqueID)

	fmt.Printf("Fetching posts (max count: %d)...\n", resolved.Count)
	var allPosts []douyin.AwemeDetail
	var cursor int64 = 0
	for {
		resp, err := cr.GetUserPosts(secUID, cursor, resolved.Count)
		if err != nil {
			fmt.Printf("Error fetching user posts: %v\n", err)
			os.Exit(1)
		}

		if len(resp.AwemeList) == 0 {
			break
		}

		allPosts = append(allPosts, resp.AwemeList...)
		fmt.Printf("Fetched %d posts so far...\n", len(allPosts))

		if resp.HasMore == 0 || resp.MaxCursor == 0 {
			break
		}
		cursor = resp.MaxCursor
		time.Sleep(1 * time.Second) // rate limit protection
	}

	if len(allPosts) == 0 {
		fmt.Println("No posts found for this user or request was blocked.")
		return
	}

	fmt.Printf("Found %d posts in total. Starting download concurrently (concurrency limit: %d)...\n", len(allPosts), resolved.Concurrency)
	var wg sync.WaitGroup
	sem := make(chan struct{}, resolved.Concurrency)

	for idx, post := range allPosts {
		wg.Add(1)
		go func(i int, p douyin.AwemeDetail) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Printf("[%d/%d] Starting download of post %s...\n", i+1, len(allPosts), p.AwemeID)
			if err := downloader.DownloadPost(&p, resolved.Output, resolved.Template, false); err != nil {
				fmt.Printf("   Error downloading post %s: %v\n", p.AwemeID, err)
			} else {
				fmt.Printf("[%d/%d] Completed post %s!\n", i+1, len(allPosts), p.AwemeID)
			}
		}(idx, post)
	}
	wg.Wait()

	fmt.Println("All user posts downloaded!")
}

// HandleLikes downloads all liked posts from a user profile
func HandleLikes(userArg string, cfg *config.Config, opts Options) {
	resolved := Resolve(cfg, opts)

	fmt.Println("Initializing Douyin Downloader...")
	cr, err := douyin.NewCrawler(resolved.Cookie, resolved.Proxy)
	if err != nil {
		fmt.Printf("Error creating crawler: %v\n", err)
		os.Exit(1)
	}
	if cfg.UserAgent != "" {
		cr.Client.UserAgent = cfg.UserAgent
	}

	if err := os.MkdirAll(resolved.Output, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Parsing user profile URL for likes...")
	secUID, err := client.ExtractSecUID(userArg)
	if err != nil {
		fmt.Printf("Error extracting sec_user_id: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found sec_user_id: %s. Fetching profile info...\n", secUID)

	profile, err := cr.GetUserProfile(secUID)
	if err != nil {
		fmt.Printf("Error fetching user profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("User profile: %s (ID: %s)\n", profile.User.Nickname, profile.User.UniqueID)

	fmt.Printf("Fetching liked posts (max count: %d)...\n", resolved.Count)
	var allPosts []douyin.AwemeDetail
	var cursor int64 = 0

	likesTemplate := resolved.Template
	if likesTemplate == "{nickname}/{publish_date}_{title}" {
		likesTemplate = "{nickname}/likes/{publish_date}_{title}"
	}

	for {
		resp, err := cr.GetUserLikes(secUID, cursor, resolved.Count)
		if err != nil {
			fmt.Printf("Error fetching liked posts: %v\n", err)
			os.Exit(1)
		}

		if len(resp.AwemeList) == 0 {
			break
		}

		allPosts = append(allPosts, resp.AwemeList...)
		fmt.Printf("Fetched %d liked posts so far...\n", len(allPosts))

		if resp.HasMore == 0 || resp.MaxCursor == 0 {
			break
		}
		cursor = resp.MaxCursor
		time.Sleep(1 * time.Second) // rate limit protection
	}

	if len(allPosts) == 0 {
		fmt.Println("No liked posts found for this user (they may be private or request was blocked).")
		return
	}

	fmt.Printf("Found %d liked posts in total. Starting download concurrently (concurrency limit: %d)...\n", len(allPosts), resolved.Concurrency)
	var wg sync.WaitGroup
	sem := make(chan struct{}, resolved.Concurrency)

	for idx, post := range allPosts {
		wg.Add(1)
		go func(i int, p douyin.AwemeDetail) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Printf("[%d/%d] Starting download of liked post %s...\n", i+1, len(allPosts), p.AwemeID)
			if err := downloader.DownloadPost(&p, resolved.Output, likesTemplate, false); err != nil {
				fmt.Printf("   Error downloading liked post %s: %v\n", p.AwemeID, err)
			} else {
				fmt.Printf("[%d/%d] Completed liked post %s!\n", i+1, len(allPosts), p.AwemeID)
			}
		}(idx, post)
	}
	wg.Wait()

	fmt.Println("All liked posts downloaded!")
}
