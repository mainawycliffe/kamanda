package utils

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/mainawycliffe/kamanda/firebase/auth"
)

func TestFormatTimestampToDate(t *testing.T) {
	os.Setenv("TZ", "UTC")
	type args struct {
		timestamp int64
		format    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Date Format is Correct",
			args: args{
				timestamp: int64(1587541614335),
				format:    "02/01/2006",
			},
			want: "22/04/2020",
		},
		{
			name: "Test Date Time Format is Correct",
			args: args{
				timestamp: int64(1587541614335),
				format:    "02/01/2006 15:04:05 MST",
			},
			want: "22/04/2020 07:46:54 UTC",
		},
		{
			name: "Test Time Format is Correct",
			args: args{
				timestamp: int64(1587541614335),
				format:    "15:04:05 MST",
			},
			want: "07:46:54 UTC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTimestampToDate(tt.args.timestamp, tt.args.format); got != tt.want {
				t.Errorf("FormatTimestampToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStringToActualValueType(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: `Test \"true\" to true boolean conversion`,
			args: args{
				input: "true",
			},
			want: true,
		},
		{
			name: `Test \"1\" to 1 int conversion`,
			args: args{
				input: "1",
			},
			want: 1,
		},
		{
			name: `Test \"1.5\" to 1.5 int conversion`,
			args: args{
				input: "1.5",
			},
			want: 1.5,
		},
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
		name       string
		args       args
		wantLength int
	}{
		{
			name: "Generate an 8 Char Password",
			args: args{
				passwordLength: 8,
			},
			wantLength: 8,
		},
		{
			name: "Generate an 10 Char Password",
			args: args{
				passwordLength: 10,
			},
			wantLength: 10,
		},
		{
			name: "Generate an 12 Char Password",
			args: args{
				passwordLength: 12,
			},
			wantLength: 12,
		},
		{
			name: "Generate an 16 Char Password",
			args: args{
				passwordLength: 16,
			},
			wantLength: 16,
		},
		{
			name: "Generate an 100 Char Password",
			args: args{
				passwordLength: 100,
			},
			wantLength: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PasswordGenerator(tt.args.passwordLength); len(got) != tt.wantLength {
				t.Errorf("PasswordGenerator() Length = %v, want %v", got, tt.wantLength)
			}
		})
	}
}

func TestPasswordGeneratorRandom(t *testing.T) {
	type args struct {
		passwordLength int
	}
	tests := []struct {
		name       string
		args       args
		wantLength int
	}{
		{
			name: "Generate an 8 Char Password",
			args: args{
				passwordLength: 8,
			},
			wantLength: 8,
		},
		{
			name: "Generate an 10 Char Password",
			args: args{
				passwordLength: 10,
			},
			wantLength: 10,
		},
		{
			name: "Generate an 12 Char Password",
			args: args{
				passwordLength: 12,
			},
			wantLength: 12,
		},
		{
			name: "Generate an 16 Char Password",
			args: args{
				passwordLength: 16,
			},
			wantLength: 16,
		},
		{
			name: "Generate an 100 Char Password",
			args: args{
				passwordLength: 100,
			},
			wantLength: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PasswordGenerator(tt.args.passwordLength); len(got) != tt.wantLength {
				t.Errorf("PasswordGenerator() Length = %v, want %v", got, tt.wantLength)
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
		{
			name: "Custom Claims Str to Interface 1 (String)",
			args: args{
				input: map[string]string{"hello": "world"},
			},
			want: map[string]interface{}{"hello": "world"},
		},
		{
			name: "Custom Claims Str to Interface 2 (Boolean)",
			args: args{
				input: map[string]string{"hello": "true"},
			},
			want: map[string]interface{}{"hello": true},
		},
		{
			name: "Custom Claims Str to Interface 3 (Int)",
			args: args{
				input: map[string]string{"hello": "1"},
			},
			want: map[string]interface{}{"hello": 1},
		},
		{
			name: "Custom Claims Str to Interface 4 (Float)",
			args: args{
				input: map[string]string{"hello": "1.23"},
			},
			want: map[string]interface{}{"hello": 1.23},
		},
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
		{
			name: "Hello World",
			args: args{
				format: "Hello %s",
				a:      []interface{}{"World"},
			},
			wantW: "\u001b[31mHello World\u001b[0m",
		},
		{
			name: "Test number 2",
			args: args{
				format: "This is a go %s",
				a:      []interface{}{"test"},
			},
			wantW: "\u001b[31mThis is a go test\u001b[0m",
		},
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
		{
			name: "Hello World",
			args: args{
				format: "Hello %s",
				a:      []interface{}{"World"},
			},
			wantW: "\u001b[32mHello World\u001b[0m",
		},
		{
			name: "Test number 2",
			args: args{
				format: "This is a go %s",
				a:      []interface{}{"test"},
			},
			wantW: "\u001b[32mThis is a go test\u001b[0m",
		},
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
	testUserResponse := []auth.FirebaseUser{
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
	testUserWrongResponse := []auth.FirebaseUser{
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
	var test1Users []auth.FirebaseUser
	var test2Users []auth.FirebaseUser
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantResponse []auth.FirebaseUser
	}{
		{
			name: "Test with JSON File",
			args: args{
				path:      "./../testdata/users.json",
				extension: "json",
				v:         &test1Users,
			},
			wantErr:      false,
			wantResponse: testUserResponse,
		},
		{
			name: "Test with Yaml File",
			args: args{
				path:      "./../testdata/users.yaml",
				extension: "yaml",
				v:         &test2Users,
			},
			wantErr:      false,
			wantResponse: testUserResponse,
		},
		{
			name: "No file",
			args: args{
				path:      "./../testdata/users1.yaml",
				extension: "yaml",
				v:         &test2Users,
			},
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name: "Unsupported Format",
			args: args{
				path:      "./../testdata/users.csv",
				extension: "csv",
				v:         &test2Users,
			},
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name: "Test with Yaml File (Incorrect Response)",
			args: args{
				path:      "./../testdata/users.yaml",
				extension: "yaml",
				v:         &test2Users,
			},
			wantErr:      false,
			wantResponse: testUserWrongResponse,
		},
		{
			name: "Test with Wrong JSON File",
			args: args{
				path:      "./../testdata/users.csv",
				extension: "json",
				v:         &test1Users,
			},
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name: "Test with Wrong Yaml File",
			args: args{
				path:      "./../testdata/users.csv",
				extension: "yaml",
				v:         &test2Users,
			},
			wantErr:      true,
			wantResponse: nil,
		},
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
