package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

// ProcessCustomClaimInput take in the input from cmd flags which a map of strings
// and convert it to a map of interface
func ProcessCustomClaimInput(input map[string]string) map[string]interface{} {
	customClaims := make(map[string]interface{})
	for k, v := range input {
		// @todo try and determine the value type and return it natively
		customClaims[k] = v
	}
	return customClaims
}

// StdOutError print an error message to the standard out
func StdOutError(format string, a ...interface{}) {
	m := aurora.Sprintf(aurora.Red(format), a)
	fmt.Fprintf(os.Stdout, "%s\n", m)
}

// StdOutSuccess print a success message to the standard out
func StdOutSuccess(format string, a ...interface{}) {
	m := aurora.Sprintf(aurora.Green(format), a)
	fmt.Fprintf(os.Stdout, "%s\n", m)
}

// UnmashalFormatFile read and unmashal either a json/yaml file into a struct
func UnmashalFormatFile(path string, extension string, v interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading %s: %w", path, err)
	}
	switch extension {
	case "yaml":
		err = yaml.Unmarshal(content, v)
		if err != nil {
			return fmt.Errorf("Error decoding yaml: %w", err)
		}
	case "json":
		err = json.Unmarshal(content, v)
		if err != nil {
			return fmt.Errorf("Error decoding json: %w", err)
		}
	default:
		return fmt.Errorf("Unsupported file type")
	}
	return nil
}
