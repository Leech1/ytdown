package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const userAgent = "com.google.android.apps.youtube.vr.oculus/1.65.10 (Linux; U; Android 12L; eureka-user Build/SQ3A.220605.009.A1) gzip"

// Streams the content at streamURL into a file at destPath.
// The User-Agent must match the client that resolved the URL, since
// YouTube's stream servers bind signed URLs to it.
func DownloadToFile(streamURL, destPath string) error {
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("requesting stream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("writing stream to file: %w", err)
	}

	return nil
}
