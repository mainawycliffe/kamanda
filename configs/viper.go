package configs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func UnsetViperConfig(keys ...string) error {
	configMap := viper.AllSettings()
	for _, key := range keys {
		delete(configMap, strings.ToLower(key))
	}
	encodedConfig, err := json.MarshalIndent(configMap, "", " ")
	if err != nil {
		return err
	}
	if err = viper.ReadConfig(bytes.NewReader(encodedConfig)); err != nil {
		return err
	}
	if err = viper.WriteConfig(); err != nil {
		return fmt.Errorf("Error removing configs: %w", err)
	}
	return nil
}
