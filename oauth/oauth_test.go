// Package oauth get refresh token and save for future use
package oauth

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/mainawycliffe/kamanda/configs"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Test_generateOauthStateTracker(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test returned string is based64 - 1"},
		{"Test returned string is based64 - 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateOauthStateTracker()
			if got == "" {
				t.Errorf("generateOauthStateTracker() = Response Required")
			}
			if _, err := base64.StdEncoding.DecodeString(got); err != nil {
				t.Errorf("generateOauthStateTracker() = %v is not a valid base64 string", got)
			}
		})
	}
}

func Test_getUserDataFromGoogle(t *testing.T) {
	type args struct {
		googleOauthConfig *oauth2.Config
		code              string
	}
	tests := []struct {
		name    string
		args    args
		want    *GetUserDataFromGoogleResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserDataFromGoogle(tt.args.googleOauthConfig, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserDataFromGoogle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUserDataFromGoogle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGoogleOAuthConfig(t *testing.T) {
	type args struct {
		port string
	}
	tests := []struct {
		name string
		args args
		want *oauth2.Config
	}{
		{
			"Test Google OAuth Configs",
			args{
				port: "200",
			},
			&oauth2.Config{
				RedirectURL:  "http://localhost:200",
				Scopes:       googleOAuthScopes,
				Endpoint:     google.Endpoint,
				ClientID:     viper.GetString(configs.GoogleOAuthClientIDConfigKey),
				ClientSecret: viper.GetString(configs.GoogleOAuthClientSecretConfigKey),
			},
		},
		{
			"Test Google OAuth Configs without Port",
			args{
				port: "",
			},
			&oauth2.Config{
				RedirectURL:  noPortURL,
				Scopes:       googleOAuthScopes,
				Endpoint:     google.Endpoint,
				ClientID:     viper.GetString(configs.GoogleOAuthClientIDConfigKey),
				ClientSecret: viper.GetString(configs.GoogleOAuthClientSecretConfigKey),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGoogleOAuthConfig(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGoogleOAuthConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
