package extractor

import (
	"errors"
	"net/url"
)

var ErrInvalidSignatureCipher = errors.New("invalid signature cipher format")

// Holds the parsed components of a Format's SignatureCipher field.
type cipherParts struct {
	Signature string // obfuscated signature, needs decoding
	SigParam  string // query param name the decoded signature should use, e.g. "sig"
	StreamURL string // base stream URL, missing its signature
}

// Parses the raw SignatureCipher query string into its components.
func parseSignatureCipher(cipher string) (*cipherParts, error) {
	values, err := url.ParseQuery(cipher)
	if err != nil {
		return nil, err
	}

	streamURL := values.Get("url")
	signature := values.Get("s")
	sigParam := values.Get("sp")

	if streamURL == "" || signature == "" {
		return nil, ErrInvalidSignatureCipher
	}

	if sigParam == "" {
		sigParam = "sig"
	}

	return &cipherParts{
		Signature: signature,
		SigParam:  sigParam,
		StreamURL: streamURL,
	}, nil
}
