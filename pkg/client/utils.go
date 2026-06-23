package client

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	secUIDPattern1 = regexp.MustCompile(`sec_uid=([^&]*)`)
	secUIDPattern2 = regexp.MustCompile(`user/([^/?]*)`)

	awemeIDPattern1 = regexp.MustCompile(`video/([^/?]*)`)
	awemeIDPattern2 = regexp.MustCompile(`[?&]vid=(\d+)`)
	awemeIDPattern3 = regexp.MustCompile(`note/([^/?]*)`)
	awemeIDPattern4 = regexp.MustCompile(`modal_id=([0-9]+)`)

	urlExtractorPattern = regexp.MustCompile(`https?://[^\s]+`)
)

// ExtractURL extracts the first valid HTTP/HTTPS URL from a string (e.g. from shared text)
func ExtractURL(text string) (string, error) {
	match := urlExtractorPattern.FindString(text)
	if match == "" {
		return "", fmt.Errorf("no URL found in input text")
	}
	return match, nil
}

// ResolveRedirect follows HTTP redirects and returns the final destination URL
func ResolveRedirect(rawURL string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Do not follow redirects automatically if we want to trace them,
			// but here we just return nil to allow standard redirection tracking.
			return nil
		},
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Request.URL.String(), nil
}

// ExtractSecUID extracts the sec_user_id from a Douyin URL (resolving redirects if necessary)
func ExtractSecUID(rawInput string) (string, error) {
	extractedURL, err := ExtractURL(rawInput)
	if err != nil {
		return "", err
	}

	finalURL := extractedURL
	if strings.Contains(extractedURL, "v.douyin.com") {
		resolved, err := ResolveRedirect(extractedURL)
		if err != nil {
			return "", fmt.Errorf("failed to resolve redirect for %s: %w", extractedURL, err)
		}
		finalURL = resolved
	}

	// Try pattern 1
	if match := secUIDPattern1.FindStringSubmatch(finalURL); len(match) > 1 {
		return match[1], nil
	}

	// Try pattern 2
	if match := secUIDPattern2.FindStringSubmatch(finalURL); len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("could not extract sec_user_id from URL: %s", finalURL)
}

// ExtractAwemeID extracts the aweme_id (video or note ID) from a Douyin URL (resolving redirects if necessary)
func ExtractAwemeID(rawInput string) (string, error) {
	extractedURL, err := ExtractURL(rawInput)
	if err != nil {
		return "", err
	}

	finalURL := extractedURL
	if strings.Contains(extractedURL, "v.douyin.com") {
		resolved, err := ResolveRedirect(extractedURL)
		if err != nil {
			return "", fmt.Errorf("failed to resolve redirect for %s: %w", extractedURL, err)
		}
		finalURL = resolved
	}

	// Order of patterns to try matching
	patterns := []*regexp.Regexp{
		awemeIDPattern1,
		awemeIDPattern2,
		awemeIDPattern3,
		awemeIDPattern4,
	}

	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(finalURL); len(match) > 1 {
			return match[1], nil
		}
	}

	// Check if the URL path itself is just the ID (e.g. after redirects)
	u, err := url.Parse(finalURL)
	if err == nil {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			// Check if it's numeric only (most aweme_ids are)
			numericPattern := regexp.MustCompile(`^\d+$`)
			if numericPattern.MatchString(lastPart) {
				return lastPart, nil
			}
		}
	}

	return "", fmt.Errorf("could not extract aweme_id from URL: %s", finalURL)
}
