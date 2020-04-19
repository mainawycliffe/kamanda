package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// FormatResults Takes any results and formats it in the requested way
func FormatResults(results interface{}, format string) (output []byte, err error) {
	if format == "json" {
		json, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("Error marsalling json: %s", err.Error())
		}
		return json, nil
	}
	if format == "yaml" {
		yaml, err := yaml.Marshal(results)
		if err != nil {
			return nil, fmt.Errorf("Error marsalling yaml: %s", err.Error())
		}
		return yaml, nil
	}
	// unrecognized format
	return nil, errors.New("Unknown Format")
}
