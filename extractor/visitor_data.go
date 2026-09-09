package extractor

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var ErrVisitorDataNotFound = errors.New("visitorData not found in homepage")

// Fetches YouTube's homepage and extracts a valid
// visitorData token from the embedded ytcfg configuration. This token
// is required by the InnerTube API — requests without it are rejected
// with a "Precondition check failed" error.
func fetchVisitorData() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.youtube.com", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:154.0) Gecko/20100101 Firefox/154.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	html := string(body)

	const marker = "\nytcfg.set("
	_, after, ok := strings.Cut(html, marker)
	if !ok {
		return "", ErrVisitorDataNotFound
	}
	rest := after

	// rest starts with a JSON object; find its end by tracking brace depth.
	depth := 0
	end := -1
	for i, ch := range rest {
		switch ch {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i + 1
			}
		}
		if end != -1 {
			break
		}
	}
	if end == -1 {
		return "", ErrVisitorDataNotFound
	}

	var cfg struct {
		InnertubeContext struct {
			Client struct {
				VisitorData string `json:"visitorData"`
			} `json:"client"`
		} `json:"INNERTUBE_CONTEXT"`
	}
	if err := json.Unmarshal([]byte(rest[:end]), &cfg); err != nil {
		return "", err
	}

	if cfg.InnertubeContext.Client.VisitorData == "" {
		return "", ErrVisitorDataNotFound
	}

	return cfg.InnertubeContext.Client.VisitorData, nil
}
