package extractor

import (
	"errors"
	"strings"
)

// ErrPlayerResponseNotFound is returned when the ytInitialPlayerResponse
// blob can't be located in the page HTML.
var ErrPlayerResponseNotFound = errors.New("ytInitialPlayerResponse not found in page")

const playerResponseMarker = "var ytInitialPlayerResponse = "

// Scans the watch page HTML and returns the raw
// JSON text of the ytInitialPlayerResponse object (unparsed, as a string).
func extractPlayerResponseJSON(html string) (string, error) {
	start := strings.Index(html, playerResponseMarker)
	if start == -1 {
		return "", ErrPlayerResponseNotFound
	}
	start += len(playerResponseMarker)

	// Finding closing brace by tracking nesting depth
	depth := 0
	for i := start; i < len(html); i++ {
		switch html[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return html[start : i+1], nil
			}
		}
	}

	return "", ErrPlayerResponseNotFound
}
