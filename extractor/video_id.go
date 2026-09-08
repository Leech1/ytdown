package extractor

import (
	"errors"
	"net/url"
)

var ErrInvalidVideoURL = errors.New("could not extract video ID from URL")

// Pulls the video ID out of a standard YouTube watch URL,
// e.g. "https://www.youtube.com/watch?v=abc123" -> "abc123".
func ExtractVideoID(videoURL string) (string, error) {
	parsed, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	id := parsed.Query().Get("v")
	if id == "" {
		return "", ErrInvalidVideoURL
	}

	return id, nil
}
