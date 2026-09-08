package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// innertubeAPIKey is YouTube's public InnerTube API key, embedded in
// every web client and safe to hardcode (it's not a secret, just an
// identifier — YouTube's own web player ships it in plain JS).
const innertubeAPIKey = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"

const innertubePlayerURL = "https://www.youtube.com/youtubei/v1/player?key=" + innertubeAPIKey

// innertubeRequestBody mirrors the minimal JSON body YouTube's InnerTube
// API expects. We declare an Android client context, since the web
// client is increasingly restricted to SABR-only streaming and omits
// direct/cipherable format URLs.
type innertubeRequestBody struct {
	VideoID string           `json:"videoId"`
	Context innertubeContext `json:"context"`
}

type innertubeContext struct {
	Client innertubeClient `json:"client"`
}

type innertubeClient struct {
	ClientName        string `json:"clientName"`
	ClientVersion     string `json:"clientVersion"`
	AndroidSDKVersion int    `json:"androidSdkVersion,omitempty"`
}

// Requests the player response directly from
// YouTube's InnerTube API using an Android client context, which
// returns real, cipherable format URLs instead of the web client's
// SABR-restricted formats.
func FetchPlayerResponseViaAPI(videoID string) (*PlayerResponse, error) {
	reqBody := innertubeRequestBody{
		VideoID: videoID,
		Context: innertubeContext{
			Client: innertubeClient{
				ClientName:        "ANDROID",
				ClientVersion:     "19.09.37",
				AndroidSDKVersion: 30,
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, innertubePlayerURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting player API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var pr PlayerResponse
	if err := json.Unmarshal(respBytes, &pr); err != nil {
		return nil, fmt.Errorf("parsing player response JSON: %w", err)
	}

	return &pr, nil
}
