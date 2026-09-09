package extractor

import (
	"fmt"

	"github.com/dop251/goja"
)

// Runs YouTube's own signature-scrambling function
// (extracted from the player JS) inside a JS engine, and returns the
// real, usable signature.
func decipherSignature(playerJS, obfuscatedSig string) (string, error) {
	funcName, err := extractDecipherFuncName(playerJS)
	if err != nil {
		return "", fmt.Errorf("finding decipher function name: %w", err)
	}

	funcSource, err := extractFunctionSource(playerJS, funcName)
	if err != nil {
		return "", fmt.Errorf("extracting decipher function source: %w", err)
	}

	helperName, err := extractHelperObjName(funcSource)
	if err != nil {
		return "", fmt.Errorf("finding helper object name: %w", err)
	}

	helperSource, err := extractHelperObjSource(playerJS, helperName)
	if err != nil {
		return "", fmt.Errorf("extracting helper object source: %w", err)
	}

	vm := goja.New()

	// Define the helper object first, then the decipher function that
	// depends on it, then call the function with our obfuscated signature.
	script := fmt.Sprintf(`
		%s
		var %s = %s;
		%s("%s");
	`, helperSource, funcName, funcSource[len(funcName)+1:], funcName, obfuscatedSig)

	result, err := vm.RunString(script)
	if err != nil {
		return "", fmt.Errorf("running decipher script in JS engine: %w", err)
	}

	return result.String(), nil
}
