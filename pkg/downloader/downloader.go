package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thaunq95/douyin-go/pkg/douyin"
)

// SanitizeSegment removes invalid directory/file name characters
func SanitizeSegment(name string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r"}
	for _, c := range invalidChars {
		name = strings.ReplaceAll(name, c, "_")
	}
	return strings.TrimSpace(name)
}

// FormatPath formats a template with nickname, publish date, title, and aweme id.
// It splits the template by slashes to build correct nested folder paths on all OS.
func FormatPath(template string, detail *douyin.AwemeDetail) string {
	publishTime := time.Unix(detail.CreateTime, 0)
	publishDate := publishTime.Format("2006-01-02")

	nickname := SanitizeSegment(detail.Author.Nickname)
	if nickname == "" {
		nickname = "unknown_user"
	}

	title := detail.Desc
	runes := []rune(title)
	if len(runes) > 50 {
		title = string(runes[:50])
	}
	title = SanitizeSegment(title)
	if title == "" {
		title = detail.AwemeID
	}

	awemeID := SanitizeSegment(detail.AwemeID)

	// Replace separators to standard forward slash
	normalizedTemplate := strings.ReplaceAll(template, "\\", "/")
	parts := strings.Split(normalizedTemplate, "/")

	var formattedParts []string
	for _, part := range parts {
		if part == "" {
			continue
		}
		r := strings.NewReplacer(
			"{nickname}", nickname,
			"{title}", title,
			"{publish_date}", publishDate,
			"{aweme_id}", awemeID,
		)
		formatted := r.Replace(part)
		formatted = SanitizeSegment(formatted)
		if formatted != "" {
			formattedParts = append(formattedParts, formatted)
		}
	}

	return filepath.Join(formattedParts...)
}

// DownloadFile downloads a URL and writes it to a destination path, tracking progress
func DownloadFile(urlStr, destPath string, showProgress bool) error {
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	tmpFile := destPath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer func() {
		out.Close()
		if _, err := os.Stat(tmpFile); err == nil {
			os.Remove(tmpFile)
		}
	}()

	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status: %d", resp.StatusCode)
	}

	totalSize := resp.ContentLength
	var downloaded int64
	buffer := make([]byte, 32*1024)

	lastUpdate := time.Now()

	for {
		n, readErr := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)

			if showProgress && time.Since(lastUpdate) > 500*time.Millisecond {
				if totalSize > 0 {
					percentage := float64(downloaded) / float64(totalSize) * 100
					fmt.Printf("\rDownloading: %.2f%% (%.2f MB / %.2f MB)...", percentage, float64(downloaded)/1024/1024, float64(totalSize)/1024/1024)
				} else {
					fmt.Printf("\rDownloading: %.2f MB...", float64(downloaded)/1024/1024)
				}
				lastUpdate = time.Now()
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	if showProgress {
		if totalSize > 0 {
			fmt.Printf("\rDownloading: 100.00%% (%.2f MB / %.2f MB) - Done!\n", float64(totalSize)/1024/1024, float64(totalSize)/1024/1024)
		} else {
			fmt.Printf("\rDownloading: %.2f MB - Done!\n", float64(downloaded)/1024/1024)
		}
	}

	out.Close()
	return os.Rename(tmpFile, destPath)
}

// DownloadPost downloads the content of a post (cover, video, audio/mp3, or slideshow images) concurrently
func DownloadPost(detail *douyin.AwemeDetail, outputDir, nameTemplate string, showProgress bool) error {
	if nameTemplate == "" {
		nameTemplate = "{nickname}/{publish_date}_{title}"
	}

	targetSubDir := FormatPath(nameTemplate, detail)
	targetDir := filepath.Join(outputDir, targetSubDir)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	if showProgress {
		fmt.Printf("Saving post assets to: %s\n", targetDir)
	}

	// Save metadata to data.json
	jsonData, err := json.MarshalIndent(detail, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal detail to JSON: %w", err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "data.json"), jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data.json: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// 1. Download Cover Image
	var coverURL string
	if len(detail.Video.Cover.URLList) > 0 {
		coverURL = detail.Video.Cover.URLList[0]
	} else if len(detail.Video.OriginCover.URLList) > 0 {
		coverURL = detail.Video.OriginCover.URLList[0]
	}
	if coverURL != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if showProgress {
				fmt.Println("-> Downloading cover image...")
			}
			// Cover download does not need a progress bar
			if err := DownloadFile(coverURL, filepath.Join(targetDir, "cover.jpg"), false); err != nil {
				errChan <- fmt.Errorf("failed to download cover image: %w", err)
			}
		}()
	}

	// 2. Download MP3 Audio
	if len(detail.Music.PlayURL.URLList) > 0 && detail.Music.PlayURL.URLList[0] != "" {
		audioURL := detail.Music.PlayURL.URLList[0]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if showProgress {
				fmt.Println("-> Downloading audio (mp3)...")
			}
			// Audio download does not need a progress bar
			if err := DownloadFile(audioURL, filepath.Join(targetDir, "audio.mp3"), false); err != nil {
				errChan <- fmt.Errorf("failed to download audio: %w", err)
			}
		}()
	}

	// 3. Download Video or Slideshow Images
	// aweme_type: 68 is image list (slideshow)
	if detail.AwemeType == 68 || len(detail.Images) > 0 {
		if showProgress {
			fmt.Printf("-> Post %s is an image slideshow. Downloading %d images...\n", detail.AwemeID, len(detail.Images))
		}
		var imgWg sync.WaitGroup
		sem := make(chan struct{}, 5) // limit slideshow image concurrency to 5
		var downloadedCount int32
		totalImages := len(detail.Images)

		for idx, img := range detail.Images {
			if len(img.URLList) == 0 {
				continue
			}
			imgURL := img.URLList[0]
			destPath := filepath.Join(targetDir, fmt.Sprintf("image_%d.jpg", idx+1))

			imgWg.Add(1)
			go func(i int, url string, path string) {
				defer imgWg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				// Individual slideshow images do not need individual progress bars
				if err := DownloadFile(url, path, false); err != nil {
					errChan <- fmt.Errorf("failed to download image %d: %w", i+1, err)
				} else {
					newCount := atomic.AddInt32(&downloadedCount, 1)
					if showProgress {
						fmt.Printf("\rDownloading slideshow images: %d/%d...", newCount, totalImages)
					}
				}
			}(idx, imgURL, destPath)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			imgWg.Wait()
			if showProgress && totalImages > 0 {
				fmt.Printf("\rDownloading slideshow images: %d/%d - Done!\n", totalImages, totalImages)
			}
		}()
	} else {
		if len(detail.Video.PlayAddr.URLList) > 0 {
			videoURL := detail.Video.PlayAddr.URLList[0]
			wg.Add(1)
			go func() {
				defer wg.Done()
				if showProgress {
					fmt.Println("-> Downloading video...")
				}
				if err := DownloadFile(videoURL, filepath.Join(targetDir, "video.mp4"), showProgress); err != nil {
					errChan <- fmt.Errorf("failed to download video: %w", err)
				}
			}()
		} else {
			if showProgress {
				fmt.Printf("   Warning: no play addresses found for video post %s\n", detail.AwemeID)
			}
		}
	}

	// Wait for all downloads to finish
	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
