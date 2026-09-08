package extractor

import (
	"errors"
	"regexp"
	"strconv"
)

var ErrDecipherFuncNotFound = errors.New("decipher function not found in player JS")

// Finds a function definition matching the pattern
// YouTube uses for signature deciphering: a function taking one argument
// that immediately does `a=a.split("")`. The function name is captured
// since it's randomized/minified on every YouTube deploy.
//
// Example match target:
//
//	fnA=function(a){a=a.split("");Xyz.aB(a,3);Xyz.reverse(a);return a.join("")}
var decipherFuncPattern = regexp.MustCompile(
	`([a-zA-Z0-9$]{2,})=function\(a\)\{a=a\.split\(""\)(.*?)return a\.join\(""\)\}`,
)

// Finds individual operation calls within a decipher
// function body, e.g. "Xyz.aB(a,3)" -> helper object "Xyz", method "aB",
// numeric argument "3".
var opCallPattern = regexp.MustCompile(`([a-zA-Z0-9$]+)\.([a-zA-Z0-9$]+)\(a,(\d+)\)`)

// Represents one operation call found in the decipher function body,
// before we know what the operation actually does.
type rawOp struct {
	HelperObj string
	Method    string
	Arg       int
}

// Finds the decipher function in the player JS and
// returns the ordered sequence of operation calls it makes.
func extractRawOpSequence(playerJS string) ([]rawOp, error) {
	match := decipherFuncPattern.FindStringSubmatch(playerJS)
	if len(match) < 3 {
		return nil, ErrDecipherFuncNotFound
	}

	body := match[2]

	calls := opCallPattern.FindAllStringSubmatch(body, -1)
	if len(calls) == 0 {
		return nil, ErrDecipherFuncNotFound
	}

	ops := make([]rawOp, 0, len(calls))
	for _, c := range calls {
		arg, err := strconv.Atoi(c[3])
		if err != nil {
			return nil, err
		}
		ops = append(ops, rawOp{
			HelperObj: c[1],
			Method:    c[2],
			Arg:       arg,
		})
	}

	return ops, nil
}
