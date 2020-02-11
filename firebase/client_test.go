package firebase

import (
	"testing"
)

func TestFirebase_setProjectID(t *testing.T) {
	type args struct {
		projectAlias              string
		firebaseProjectConfigFile string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantProjectId string
	}{
		{
			"Test Case",
			args{
				projectAlias:              "default",
				firebaseProjectConfigFile: "./../testdata/.firebaserc",
			},
			false,
			"kamanda-test-project",
		},
		{
			"Test Case 2",
			args{
				firebaseProjectConfigFile: "./../testdata/.firebaserc",
				projectAlias:              "default2",
			},
			true,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Firebase{}
			err := f.setProjectID(tt.args.projectAlias, tt.args.firebaseProjectConfigFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Firebase.setProjectID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if f.projectId != tt.wantProjectId {
				t.Errorf("Firebase.setProjectID() Project Id = %v, wanted %v", f.projectId, tt.wantProjectId)
			}
		})
	}
}
