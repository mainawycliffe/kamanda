// Package oauth get refresh token and save for future use
package firebase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mainawycliffe/kamanda/configs"
	"github.com/spf13/viper"
)

func setViperConfigs(setViperConfigs bool) {
	viper.SetConfigFile("./../testdata/.viper.yaml")
	viper.Reset()
	if setViperConfigs {
		viper.Set(configs.GoogleOAuthClientIDConfigKey, "ClientID")
		viper.Set(configs.GoogleOAuthClientSecretConfigKey, "ClientSecret")
		viper.Set(configs.FirebaseRefreshTokenViperConfigKey, "RefreshToken")
	}
}

func Test_constructToken(t *testing.T) {

	want, _ := json.Marshal(RefreshToken{
		ClientID:     "ClientID",
		ClientSecret: "ClientSecret",
		RefreshToken: "RefreshToken",
		Type:         "authorized_user",
	})

	tests := []struct {
		name            string
		want            []byte
		wantErr         bool
		setViperConfigs bool
	}{
		{"Test Empty Token Retrieval", nil, true, false},
		{"Test Token Retrieval", want, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setViperConfigs(tt.setViperConfigs)
			got, err := constructToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("constructToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("constructToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken_validate(t *testing.T) {
	type fields struct {
		ClientID     string
		ClientSecret string
		RefreshToken string
		Type         string
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		errResponse error
	}{
		{
			"Test No Error",
			fields{
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				RefreshToken: "RefreshToken",
				Type:         "authorized_user",
			},
			false,
			nil,
		},
		{
			"Test No Client ID Error",
			fields{
				ClientID:     "",
				ClientSecret: "ClientSecret",
				RefreshToken: "RefreshToken",
				Type:         "authorized_user",
			},
			true,
			fmt.Errorf("client id can not be empty"),
		},
		{
			"Test No Client Secret Error",
			fields{
				ClientID:     "ClientID",
				ClientSecret: "",
				RefreshToken: "RefreshToken",
				Type:         "authorized_user",
			},
			true,
			fmt.Errorf("client secret can not be empty"),
		},
		{
			"Test No Refresh Token Error",
			fields{
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				RefreshToken: "",
				Type:         "authorized_user",
			},
			true,
			fmt.Errorf("refresh token can not be empty"),
		},
		{
			"Test No Type Error",
			fields{
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				RefreshToken: "RefreshToken",
				Type:         "",
			},
			true,
			fmt.Errorf("token type is not valid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RefreshToken{
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				RefreshToken: tt.fields.RefreshToken,
				Type:         tt.fields.Type,
			}
			err := r.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.errResponse == nil {
				return
			}
			if err == nil || tt.errResponse == nil {
				t.Errorf("RefreshToken.validate() error = %v, errResponse = %v", err, tt.errResponse)
				return
			}
			if err.Error() != tt.errResponse.Error() {
				t.Errorf("RefreshToken.validate() error = %v, errResponse = %v", err, tt.errResponse)
			}
		})
	}
}
