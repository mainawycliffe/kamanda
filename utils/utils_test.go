package utils

import (
	"reflect"
	"testing"
)

func TestParseStringToActualValueType(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{`Test \"true\" to true boolean conversion`, args{input: "true"}, true},
		{`Test \"1\" to 1 int conversion`, args{input: "1"}, 1},
		{`Test \"1.5\" to 1.5 int conversion`, args{input: "1.5"}, 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseStringToActualValueType(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStringToActualValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordGenerator(t *testing.T) {
	type args struct {
		passwordLength int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Generate an 8 Char Password", args{passwordLength: 8}, 8},
		{"Generate an 10 Char Password", args{passwordLength: 10}, 10},
		{"Generate an 12 Char Password", args{passwordLength: 12}, 12},
		{"Generate an 16 Char Password", args{passwordLength: 16}, 16},
		{"Generate an 100 Char Password", args{passwordLength: 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PasswordGenerator(tt.args.passwordLength); len(got) != tt.want {
				t.Errorf("PasswordGenerator() Length = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessCustomClaimInput(t *testing.T) {
	type args struct {
		input map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"Custom Claims Str to Interface 1 (String)", args{input: map[string]string{"hello": "world"}}, map[string]interface{}{"hello": "world"}},
		{"Custom Claims Str to Interface 2 (Boolean)", args{input: map[string]string{"hello": "true"}}, map[string]interface{}{"hello": true}},
		{"Custom Claims Str to Interface 2 (Int)", args{input: map[string]string{"hello": "1"}}, map[string]interface{}{"hello": 1}},
		{"Custom Claims Str to Interface 2 (Float)", args{input: map[string]string{"hello": "1.23"}}, map[string]interface{}{"hello": 1.23}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessCustomClaimInput(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessCustomClaimInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

