package firebase

import (
	"context"
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

func TestFirebase_initializeFirebaseApp(t *testing.T) {
	type args struct {
		ctx          context.Context
		projectAlias string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantAppToBeNil bool
		wantProjectId  string
	}{
		{
			"Test Firebase App Init",
			args{
				ctx:          context.Background(),
				projectAlias: "default",
			},
			false,
			false,
			"kamanda-test-project",
		},
		{
			"Test Firebase App Init (No Default Alias)",
			args{
				ctx:          context.Background(),
				projectAlias: "",
			},
			false,
			false,
			"kamanda-test-project",
		},
		{
			"Test Firebase App Init (No existent Default Alias)",
			args{
				ctx:          context.Background(),
				projectAlias: "helloworld",
			},
			true,
			true,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setViperConfigs(true)
			f := &Firebase{}
			if err := f.initializeFirebaseApp(tt.args.ctx, tt.args.projectAlias); (err != nil) != tt.wantErr {
				t.Errorf("Firebase.initializeFirebaseApp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (f.App == nil) != tt.wantAppToBeNil {
				t.Errorf("Firebase.initializeFirebaseApp() Firebase.App is nil = %v, want Firebase.App to be nil = %v", (f.App == nil), tt.wantAppToBeNil)
			}
			if f.projectId != tt.wantProjectId {
				t.Errorf("Firebase.initializeFirebaseApp() Firebase.projectId = %v, want = %v", f.projectId, tt.wantProjectId)
			}
		})
	}
}
