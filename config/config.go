// Package config save configurations to file using github.com/spf13/viper
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

const refreshTokenFilePath = ".kamanda/refresh_token.json"

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
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("Error reading home dir: %w", err)
	}
	jsonToken, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		return fmt.Errorf("Error converting to json: %w", err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", home, refreshTokenFilePath), jsonToken, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("Error saving configs: %w", err)
	}
	return nil
}
