// Package oauth get refresh token and save for future use
package oauth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type RefreshToken struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

func (r RefreshToken) validate() error {
	if r.ClientID == "" {
		return fmt.Errorf("client id can not be empty")
	}
	if r.ClientSecret == "" {
		return fmt.Errorf("client secret can not be empty")
	}
	if r.RefreshToken == "" {
		return fmt.Errorf("refresh token can not be empty")
	}
	if r.Type != "authorized_user" {
		return fmt.Errorf("token type is not valid")
	}
	return nil
}

// SaveRefreshToken saves refresh token to file
func SaveRefreshToken(token RefreshToken) error {
	if err := token.validate(); err != nil {
		return fmt.Errorf("Error validating the refresh token: %w", err)
	}
	jsonToken, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		return fmt.Errorf("Error converting to json: %w", err)
	}
	refreshTokenFilePath := viper.GetString("refreshTokenFilePath")
	_, err = os.Stat(refreshTokenFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Error saving configs: %w", err)
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(viper.GetString("configDirPath"), os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating refresh token dir: %w", err)
		}
	}
	file, err := os.OpenFile(refreshTokenFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error saving refresh token: %w", err)
	}
	defer file.Close()
	_, err = file.Write(jsonToken)
	if err != nil {
		return fmt.Errorf("Error saving refresh token: %w", err)
	}
	return nil
}
