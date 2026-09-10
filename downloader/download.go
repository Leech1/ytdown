package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const userAgent = "com.google.android.apps.youtube.vr.oculus/1.65.10 (Linux; U; Android 12L; eureka-user Build/SQ3A.220605.009.A1) gzip"

// chunkSize is the size of each bounded range request. YouTube's video
// servers reject fully open-ended range requests (bytes=0-) but accept
// bounded ones, so we fetch the file in fixed-size chunks and append them.
const chunkSize = 10 * 1024 * 1024 // 10MB

// DownloadToFile streams the content at streamURL into a file at destPath,
// fetching it in bounded range chunks.
func DownloadToFile(streamURL, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	var start int64 = 0

	for {
		end := start + chunkSize - 1

		req, err := http.NewRequest(http.MethodGet, streamURL, nil)
		if err != nil {
			return fmt.Errorf("building request: %w", err)
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("requesting chunk at offset %d: %w", start, err)
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			resp.Body.Close()
			return fmt.Errorf("unexpected status code %d at offset %d", resp.StatusCode, start)
		}

		n, err := io.Copy(out, resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("writing chunk at offset %d: %w", start, err)
		}

		// If we got fewer bytes than requested, we've reached the end of the file.
		if n < chunkSize {
			break
		}

		start += chunkSize
	}

	return nil
}
