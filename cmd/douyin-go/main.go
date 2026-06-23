package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thaunq95/douyin-go/internal/app/handlers"
	"github.com/thaunq95/douyin-go/pkg/config"
)

var (
	cookieFlag      string
	outputFlag      string
	templateFlag    string
	countFlag       int
	proxyFlag       string
	concurrencyFlag int
)

func main() {
	// 1. Load config file settings
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: failed to load config.yaml, using defaults: %v\n", err)
	}

	// 2. Define Cobra commands
	var rootCmd = &cobra.Command{
		Use:   "douyin-go",
		Short: "Douyin Downloader in Go",
		Long:  `A powerful, concurrent Douyin downloader to download videos, audios, covers, and slideshows.`,
	}

	// Define persistent flags
	rootCmd.PersistentFlags().StringVarP(&cookieFlag, "cookie", "c", "", "Custom Douyin session cookies")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "Output directory for downloaded files")
	rootCmd.PersistentFlags().StringVarP(&templateFlag, "template", "t", "", "Filename naming template")
	rootCmd.PersistentFlags().StringVarP(&proxyFlag, "proxy", "p", "", "Proxy URL (http/https/socks5, e.g. http://127.0.0.1:8080)")
	rootCmd.PersistentFlags().IntVarP(&concurrencyFlag, "concurrency", "j", 0, "Number of concurrent downloads allowed")
	rootCmd.PersistentFlags().IntVarP(&countFlag, "count", "n", 0, "Number of posts to fetch when downloading user posts/likes")

	// Helper to extract flags into handlers.Options struct
	getOptions := func() handlers.Options {
		return handlers.Options{
			Cookie:      cookieFlag,
			Output:      outputFlag,
			Template:    templateFlag,
			Proxy:       proxyFlag,
			Concurrency: concurrencyFlag,
			Count:       countFlag,
		}
	}

	// Command to download a single video/post
	var videoCmd = &cobra.Command{
		Use:     "video [url/sharing text]",
		Aliases: []string{"download"},
		Short:   "Download a single video or post by sharing link/URL",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.HandleVideo(args[0], cfg, getOptions())
		},
	}

	// Command to download all user posts
	var userCmd = &cobra.Command{
		Use:   "user [profile_url/sharing text]",
		Short: "Download all published posts from a user profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.HandleUser(args[0], cfg, getOptions())
		},
	}

	// Command to download liked posts
	var likesCmd = &cobra.Command{
		Use:   "likes [profile_url/sharing text]",
		Short: "Download all liked posts from a user profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.HandleLikes(args[0], cfg, getOptions())
		},
	}

	// 3. Register subcommands and run
	rootCmd.AddCommand(videoCmd, userCmd, likesCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
