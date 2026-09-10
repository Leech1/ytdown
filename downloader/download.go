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
const chunkDelay = 500 * time.Millisecond

func DownloadToFile(streamURL, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	var start int64 = 0
	first := true

	for {
		if !first {
			time.Sleep(chunkDelay)
		}
		first = false

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

		if n < chunkSize {
			break
		}

		start += chunkSize
	}

	return nil
}
