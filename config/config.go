// Package config save configurations to file using github.com/spf13/viper
package config

type RefreshToken struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

func SaveRefreshToken(token RefreshToken) {
	panic("not implemented")
}
