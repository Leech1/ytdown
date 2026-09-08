package extractor

import "fmt"

// Returns a real, playable URL for the given format.
// If the format already has a direct URL, it's returned as-is. Otherwise,
// its SignatureCipher is parsed and deciphered using the player JS.
func ResolveFormatURL(format Format, playerJS string) (string, error) {
	if format.URL != "" {
		return format.URL, nil
	}

	if format.SignatureCipher == "" {
		return "", fmt.Errorf("format has neither URL nor SignatureCipher")
	}

	parts, err := parseSignatureCipher(format.SignatureCipher)
	if err != nil {
		return "", fmt.Errorf("parsing signature cipher: %w", err)
	}

	sig, err := decipherSignature(playerJS, parts.Signature)
	if err != nil {
		return "", fmt.Errorf("deciphering signature: %w", err)
	}

	return parts.StreamURL + "&" + parts.SigParam + "=" + sig, nil
}
