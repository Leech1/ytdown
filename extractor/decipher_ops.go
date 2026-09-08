package extractor

import (
	"errors"
	"regexp"
)

var ErrHelperObjNotFound = errors.New("helper object not found in player JS")

type opType int

const (
	opUnknown opType = iota
	opSwap
	opSplice
	opReverse
)

// Finds the helper object's full body given its name,
// e.g. for "Xyz" it matches: var Xyz={ ... };
func helperObjPattern(name string) *regexp.Regexp {
	return regexp.MustCompile(`(?:var\s+)?` + regexp.QuoteMeta(name) + `=\{(.*?)\};`)
}

// Matches a method body that swaps the first element with
// the element at index (b % length).
var swapPattern = regexp.MustCompile(`([a-zA-Z0-9$]+):function\(a,b\)\{var c=a\[0\];a\[0\]=a\[b%a\.length\];a\[b%a\.length\]=c\}`)

// Matches a method body that removes the first b elements.
var splicePattern = regexp.MustCompile(`([a-zA-Z0-9$]+):function\(a,b\)\{a\.splice\(0,b\)\}`)

// Matches a method body that reverses the array.
var reversePattern = regexp.MustCompile(`([a-zA-Z0-9$]+):function\(a\)\{a\.reverse\(\)\}`)

// Finds the named helper object in the player JS
// and returns a map from method name to the operation it performs.
func classifyHelperMethods(playerJS, helperObjName string) (map[string]opType, error) {
	match := helperObjPattern(helperObjName).FindStringSubmatch(playerJS)
	if len(match) < 2 {
		return nil, ErrHelperObjNotFound
	}
	body := match[1]

	classified := make(map[string]opType)

	if m := swapPattern.FindStringSubmatch(body); len(m) == 2 {
		classified[m[1]] = opSwap
	}
	if m := splicePattern.FindStringSubmatch(body); len(m) == 2 {
		classified[m[1]] = opSplice
	}
	if m := reversePattern.FindStringSubmatch(body); len(m) == 2 {
		classified[m[1]] = opReverse
	}

	return classified, nil
}
