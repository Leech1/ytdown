package extractor

import (
	"encoding/json"
	"fmt"
)

// Fetches the watch page and returns the raw ytInitialPlayerResponse JSON string.
func GetPlayerResponseJSON(videoURL string) (string, error) {
	html, err := fetchWatchPage(videoURL)
	if err != nil {
		return "", err
	}

	return extractPlayerResponseJSON(html)
}

// Fetches the watch page for videoID and returns the
// parsed PlayerResponse, including available stream formats.
func GetPlayerResponse(videoURL string) (*PlayerResponse, error) {
	raw, err := GetPlayerResponseJSON(videoURL)
	if err != nil {
		return nil, err
	}

	var pr PlayerResponse
	if err := json.Unmarshal([]byte(raw), &pr); err != nil {
		return nil, fmt.Errorf("parsing player response JSON: %w", err)
	}

	return &pr, nil
}

// Fetches the watch page for videoURL, locates the player JS
// file referenced in it, and returns the JS file's full source.
func GetPlayerJS(videoURL string) (string, error) {
	html, err := fetchWatchPage(videoURL)
	if err != nil {
		return "", fmt.Errorf("fetching watch page: %w", err)
	}

	jsURL, err := extractPlayerJSURL(html)
	if err != nil {
		return "", fmt.Errorf("finding player JS URL: %w", err)
	}

	return fetchPlayerJS(jsURL)
}
