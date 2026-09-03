package extractor

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var ErrPlayerJSNotFound = errors.New("player JS URL not found in page")

// Matches the base.js path YouTube references in the
// watch page, e.g. "/s/player/64d19a4b/player_ias.vflset/en_US/base.js"
var playerJSPattern = regexp.MustCompile(`"jsUrl":"([^"]+)"`)

// Finds the relative URL of the player JS file
// referenced in the watch page HTML, and returns it as an absolute URL.
func extractPlayerJSURL(html string) (string, error) {
	matches := playerJSPattern.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", ErrPlayerJSNotFound
	}

	return "https://www.youtube.com" + matches[1], nil
}

// Downloads the raw contents of the player JS file at the
// given absolute URL (as returned by extractPlayerJSURL).
func fetchPlayerJS(jsURL string) (string, error) {
	resp, err := http.Get(jsURL)
	if err != nil {
		return "", fmt.Errorf("requesting player JS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading player JS body: %w", err)
	}

	return string(body), nil
}
