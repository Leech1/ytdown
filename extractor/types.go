package extractor

type PlayerResponse struct {
	StreamingData StreamingData `json:"streamingData"`
	VideoDetails  VideoDetails  `json:"videoDetails"`
}

type VideoDetails struct {
	Title     string `json:"title"`
	VideoID   string `json:"videoId"`
	LengthS   string `json:"lengthSeconds"`
	IsLiveNow bool   `json:"isLiveContent"`
}

type StreamingData struct {
	Formats []Format `json:"formats"`

	// AdaptiveFormats are video-only or audio-only streams, available in
	// higher resolutions. Not used yet since we're targeting 720p progressive.
	AdaptiveFormats []Format `json:"adaptiveFormats"`
}

type Format struct {
	Itag         int    `json:"itag"`
	URL          string `json:"url"`
	MimeType     string `json:"mimeType"`
	Quality      string `json:"quality"`
	QualityLabel string `json:"qualityLabel"`

	// SignatureCipher is set instead of URL when the stream URL is
	// cipher-protected and needs to be deobfuscated before use.
	SignatureCipher string `json:"signatureCipher"`
}
