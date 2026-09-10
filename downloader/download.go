package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const userAgent = "com.google.android.apps.youtube.vr.oculus/1.65.10 (Linux; U; Android 12L; eureka-user Build/SQ3A.220605.009.A1) gzip"
const chunkSize = 10 * 1024 * 1024 // 10MB
const maxRetries = 5
const retryBaseDelay = 1 * time.Second

// Streams the content at streamURL into a file at destPath,
// fetching it in bounded range chunks. YouTube's video servers
// occasionally reject a chunk request with 403 for reasons that aren't
// fully deterministic (likely anti-bulk-download heuristics), so each
// chunk is retried with exponential backoff before giving up.
func DownloadToFile(streamURL, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	var start int64 = 0

	for {
		n, err := fetchChunkWithRetry(out, streamURL, start)
		if err != nil {
			return err
		}

		if n < chunkSize {
			break
		}
		start += chunkSize

		time.Sleep(3 * time.Second) // TEMP: test if pacing avoids the 403
	}

	return nil
}

// Requests a single chunk starting at offset start, retrying with exponential backoff on failure.
// Returns the number of bytes written.
func fetchChunkWithRetry(out *os.File, streamURL string, start int64) (int64, error) {
	end := start + chunkSize - 1

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryBaseDelay * time.Duration(1<<attempt) // 1s, 2s, 4s, 8s...
			time.Sleep(delay)
		}

		req, err := http.NewRequest(http.MethodGet, streamURL, nil)
		if err != nil {
			return 0, fmt.Errorf("building request: %w", err)
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("requesting chunk at offset %d: %w", start, err)
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			resp.Body.Close()
			lastErr = fmt.Errorf("unexpected status code %d at offset %d (attempt %d/%d)", resp.StatusCode, start, attempt+1, maxRetries)
			continue
		}

		n, err := io.Copy(out, resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("writing chunk at offset %d: %w", start, err)
			continue
		}

		return n, nil
	}

	return 0, fmt.Errorf("chunk at offset %d failed after %d attempts: %w", start, maxRetries, lastErr)
}
