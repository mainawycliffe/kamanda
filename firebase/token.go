// Package oauth get refresh token and save for future use
package firebase

import (
	"encoding/json"
	"fmt"

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

// constructToken create a json refresh token byte for use by a firebase client
func constructToken() ([]byte, error) {
	refreshToken := RefreshToken{
		ClientID:     viper.GetString("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: viper.GetString("GOOGLE_OAUTH_CLIENT_SECRET"),
		RefreshToken: viper.GetString("FirebaseRefreshToken"),
		Type:         "authorized_user",
	}
	if err := refreshToken.validate(); err != nil {
		return nil, fmt.Errorf("Error validating token: %w", err)
	}
	return json.Marshal(refreshToken)
}
