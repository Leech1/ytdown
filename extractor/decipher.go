package extractor

import (
	"errors"
	"regexp"
)

var ErrDecipherFuncNotFound = errors.New("decipher function not found in player JS")

// Finds the name of the function used to decipher
// signatures, by matching the call site where it's invoked, e.g.:
//
//	a.set("s",SOME_NAME(decodeURIComponent(b)))
//
// The function name itself is randomized per YouTube deploy, so we find
// it by context rather than by name.
var decipherFuncNamePattern = regexp.MustCompile(`\bset\("s",\s*([a-zA-Z0-9$]+)\(`)

// Finds the name of the signature decipher function within the player JS.
func extractDecipherFuncName(playerJS string) (string, error) {
	match := decipherFuncNamePattern.FindStringSubmatch(playerJS)
	if len(match) < 2 {
		return "", ErrDecipherFuncNotFound
	}
	return match[1], nil
}

// Finds the full source text of a function given
// its name, e.g. for name "fnA" it matches:
//
//	fnA=function(a){...}
//
// up to the matching closing brace, by tracking nesting depth.
func extractFunctionSource(playerJS, funcName string) (string, error) {
	marker := funcName + "=function("
	start := indexOf(playerJS, marker)
	if start == -1 {
		return "", ErrDecipherFuncNotFound
	}

	depth := 0
	bodyStarted := false
	for i := start; i < len(playerJS); i++ {
		switch playerJS[i] {
		case '{':
			depth++
			bodyStarted = true
		case '}':
			depth--
			if bodyStarted && depth == 0 {
				return playerJS[start : i+1], nil
			}
		}
	}

	return "", ErrDecipherFuncNotFound
}

// Tiny helper to avoid importing strings just for one call.
func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
