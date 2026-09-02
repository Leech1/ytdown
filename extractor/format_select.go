package extractor

import "errors"

// For the requested quality.
var ErrFormatNotFound = errors.New("no matching format found")

// YouTube's progressive (audio+video muxed) mp4 format at 720p.
const itag22 = 22

func Find720pProgressive(formats []Format) (*Format, error) {
	for _, f := range formats {
		if f.Itag == itag22 {
			return &f, nil
		}
	}

	for _, f := range formats {
		if f.QualityLabel == "720p" {
			return &f, nil
		}
	}

	return nil, ErrFormatNotFound
}
