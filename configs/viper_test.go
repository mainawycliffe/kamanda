package configs

import (
	"testing"

	"github.com/spf13/viper"
)

func TestUnsetViperConfig(t *testing.T) {
	viper.SetConfigFile("./../testdata/.viper.yaml")
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, v := range keys {
		viper.Set(v, v)
	}
	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test Case", args{keys: keys}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnsetViperConfig(tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("UnsetViperConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, v := range keys {
				if viper.GetString(v) != "" {
					t.Errorf("UnsetViperConfig() Expected viper key %s to be empty", v)
				}
			}
		})
	}
}
