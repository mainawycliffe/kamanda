package utils

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mainawycliffe/kamanda/firebase/auth"
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
		{"Custom Claims Str to Interface 3 (Int)", args{input: map[string]string{"hello": "1"}}, map[string]interface{}{"hello": 1}},
		{"Custom Claims Str to Interface 4 (Float)", args{input: map[string]string{"hello": "1.23"}}, map[string]interface{}{"hello": 1.23}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessCustomClaimInput(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessCustomClaimInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdOutError(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{"Hello World", args{format: "Hello %s", a: []interface{}{"World"}}, "\u001b[31mHello World\u001b[0m\n"},
		{"Test number 2", args{format: "This is a go %s", a: []interface{}{"test"}}, "\u001b[31mThis is a go test\u001b[0m\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			StdOutError(w, tt.args.format, tt.args.a...)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("StdOutError() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestStdOutSuccess(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{"Hello World", args{format: "Hello %s", a: []interface{}{"World"}}, "\u001b[32mHello World\u001b[0m\n"},
		{"Test number 2", args{format: "This is a go %s", a: []interface{}{"test"}}, "\u001b[32mThis is a go test\u001b[0m\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			StdOutSuccess(w, tt.args.format, tt.args.a...)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("StdOutSuccess() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestUnmarshalFormatFile(t *testing.T) {
	type args struct {
		path      string
		extension string
		v         interface{}
	}
	testUserResponse := []auth.NewUser{
		{
			DisplayName: "Some name here",
			Email:       "james@gmail.com",
			Password:    "HelloWorld",
		},
		{
			DisplayName: "Some name here",
			Email:       "ms@outlook.com",
			Password:    "HelloWorld",
		},
	}
	testUserWrongResponse := []auth.NewUser{
		{
			DisplayName: "Some name here1",
			Email:       "james@gmail.com1",
			Password:    "HelloWorld1",
		},
		{
			DisplayName: "Some name here1",
			Email:       "ms@outlook.com1",
			Password:    "HelloWorld1",
		},
	}
	var test1Users []auth.NewUser
	var test2Users []auth.NewUser
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantResponse []auth.NewUser
	}{
		{"Test with JSON File", args{path: "./../testdata/users.json", extension: "json", v: &test1Users}, false, testUserResponse},
		{"Test with Yaml File", args{path: "./../testdata/users.yaml", extension: "yaml", v: &test2Users}, false, testUserResponse},
		{"No file", args{path: "./../testdata/users1.yaml", extension: "yaml", v: &test2Users}, true, nil},
		{"Test with Yaml File (Incorrect Response)", args{path: "./../testdata/users.yaml", extension: "yaml", v: &test2Users}, false, testUserWrongResponse},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalFormatFile(tt.args.path, tt.args.extension, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFormatFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(&testUserResponse, tt.args.v) {
				t.Errorf("UnmarshalFormatFile() want = %v, got %v", testUserResponse, tt.args.v)
			}
		})
	}
}
