package extractor

import "regexp"

// Finds the name of a helper object referenced
// inside the decipher function body, e.g. for "Xyz.aB(a,3)" it captures
// "Xyz". We only need to find one such reference since the decipher
// function typically calls into a single helper object throughout.
var helperObjNamePattern = regexp.MustCompile(`;([a-zA-Z0-9$]+)\.[a-zA-Z0-9$]+\(a,`)

// Finds the name of the helper object used within
// a decipher function's source.
func extractHelperObjName(decipherFuncSource string) (string, error) {
	match := helperObjNamePattern.FindStringSubmatch(decipherFuncSource)
	if len(match) < 2 {
		return "", ErrDecipherFuncNotFound
	}
	return match[1], nil
}

// Finds the full source of a helper object given
// its name, e.g. for name "Xyz" it matches:
//
//	var Xyz={aB:function(a,b){...},cD:function(a){...}};
//
// up to the matching closing brace, by tracking nesting depth.
func extractHelperObjSource(playerJS, objName string) (string, error) {
	marker := objName + "={"
	start := indexOf(playerJS, marker)
	if start == -1 {
		return "", ErrDecipherFuncNotFound
	}

	// Move start back to include "var Xyz=" if present, so the extracted
	// text is a valid standalone JS statement when run through goja.
	varMarker := "var " + marker
	if varStart := indexOf(playerJS, varMarker); varStart != -1 && varStart < start {
		start = varStart
	}

	braceStart := indexOf(playerJS[start:], "{") + start
	depth := 0
	for i := braceStart; i < len(playerJS); i++ {
		switch playerJS[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				// include trailing semicolon if present
				end := i + 1
				if end < len(playerJS) && playerJS[end] == ';' {
					end++
				}
				return playerJS[start:end], nil
			}
		}
	}

	return "", ErrDecipherFuncNotFound
}
