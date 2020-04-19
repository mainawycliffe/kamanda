package utils

import (
	"reflect"
	"testing"
)

func TestFormatResults(t *testing.T) {
	type args struct {
		results interface{}
		format  string
	}
	results := map[string]interface{}{
		"name": "Jane Doe",
		"age":  13,
	}
	yamlOutput := []byte(`age: 13
name: Jane Doe
`)
	jsonOutput := []byte(`{
  "age": 13,
  "name": "Jane Doe"
}`)
	tests := []struct {
		name       string
		args       args
		wantOutput []byte
		wantErr    bool
	}{
		{
			name: "Test JSON Format Output",
			args: args{
				results: results,
				format:  "json",
			},
			wantOutput: jsonOutput,
			wantErr:    false,
		},
		{
			name: "Test Yaml Format Output",
			args: args{
				results: results,
				format:  "yaml",
			},
			wantOutput: yamlOutput,
			wantErr:    false,
		},
		{
			name: "Test Unsupported Format Output",
			args: args{
				results: results,
				format:  "csv",
			},
			wantOutput: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := FormatResults(tt.args.results, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("FormatResults() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
