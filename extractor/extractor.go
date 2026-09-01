package extractor

// Fetches the watch page and returns the raw ytInitialPlayerResponse JSON string.
func GetPlayerResponseJSON(videoUrl string) (string, error) {
	html, err := fetchWatchPage(videoUrl)
	if err != nil {
		return "", err
	}

	return extractPlayerResponseJSON(html)
}
