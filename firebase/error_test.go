package firebase

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want FirebaseError
	}{
		{
			"Test Error 1",
			args{
				err: fmt.Errorf("ERROR"),
			},
			FirebaseError{
				FirebaseError: fmt.Errorf("ERROR"),
			},
		},
		{
			"Test Error 2",
			args{
				err: errors.New("ERROR"),
			},
			FirebaseError{
				FirebaseError: errors.New("ERROR"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.err); got.Error() != tt.want.Error() {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}
