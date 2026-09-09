package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const innertubeAPIKey = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"

const innertubePlayerURL = "https://www.youtube.com/youtubei/v1/player"

const androidUserAgent = "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip"

type innertubeRequestBody struct {
	VideoID         string           `json:"videoId"`
	Context         innertubeContext `json:"context"`
	ContentCheckOK  bool             `json:"contentCheckOk"`
	RacyCheckOK     bool             `json:"racyCheckOk"`
	PlaybackContext playbackContext  `json:"playbackContext"`
}

type playbackContext struct {
	ContentPlaybackContext contentPlaybackContext `json:"contentPlaybackContext"`
}

type contentPlaybackContext struct {
	HTML5Preference string `json:"html5Preference"`
}

type innertubeContext struct {
	Client innertubeClient `json:"client"`
}

type innertubeClient struct {
	HL            string `json:"hl"`
	GL            string `json:"gl"`
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
	UserAgent     string `json:"userAgent"`
	TimeZone      string `json:"timeZone"`
	UTCOffset     int    `json:"utcOffsetMinutes"`
	VisitorData   string `json:"visitorData,omitempty"`
	// AndroidSDKVersion is deliberately omitted. Setting it
	// signals a fuller Android client capable of stricter verification,
	// which can trigger additional bot-check requirements.
}

// Requests the player response directly from
// YouTube's InnerTube API using an Android client context, which
// returns real, cipherable format URLs instead of the web client's
// SABR-restricted formats.
func FetchPlayerResponseViaAPI(videoID string) (*PlayerResponse, error) {
	_, offsetSeconds := time.Now().Zone()

	visitorData, err := fetchVisitorData()
	if err != nil {
		return nil, fmt.Errorf("fetching visitor data: %w", err)
	}

	reqBody := innertubeRequestBody{
		VideoID:        videoID,
		ContentCheckOK: true,
		RacyCheckOK:    true,
		PlaybackContext: playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				HTML5Preference: "HTML5_PREF_WANTS",
			},
		},
		Context: innertubeContext{
			Client: innertubeClient{
				HL:            "en",
				GL:            "US",
				ClientName:    "ANDROID_VR",
				ClientVersion: "1.65.10",
				UserAgent:     "com.google.android.apps.youtube.vr.oculus/1.65.10 (Linux; U; Android 12L; eureka-user Build/SQ3A.220605.009.A1) gzip",
				TimeZone:      "UTC",
				UTCOffset:     offsetSeconds / 60,
				VisitorData:   visitorData,
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
	req.Header.Set("User-Agent", "com.google.android.apps.youtube.vr.oculus/1.65.10 (Linux; U; Android 12L; eureka-user Build/SQ3A.220605.009.A1) gzip")
	req.AddCookie(&http.Cookie{
		Name:   "CONSENT",
		Value:  "YES+cb.20210328-17-p0.en+FX+100",
		Path:   "/",
		Domain: ".youtube.com",
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting player API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
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
