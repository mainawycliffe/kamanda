package firebase

import (
	"context"
	"encoding/json"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	s "github.com/bitfield/script"
	"google.golang.org/api/option"
)

const (
	firebaseProjectConfigFile = "./.firebaserc"
	defaultProject            = "default"
)

type FirebaseProjectConfigs struct {
	Projects map[string]string `json:"projects"`
	Targets  interface{}       `json:"-"`
}

type Firebase struct {
	App         *firebase.App
	projectId   string
	credentials []byte
}

// setProjectID use the project alias to get the firebase project id
func (f *Firebase) setProjectID(projectAlias string, firebaseProjectConfigFile string) error {
	configFileContent, err := s.File(firebaseProjectConfigFile).Bytes()
	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}
	var decodedConfigs *FirebaseProjectConfigs
	err = json.Unmarshal(configFileContent, &decodedConfigs)
	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}
	if decodedConfigs.Projects[projectAlias] == "" {
		return fmt.Errorf("Couldn't find the project alias provided!")
	}
	f.projectId = decodedConfigs.Projects[projectAlias]
	return nil
}

// initializeFirebaseApp create a new firebase app that can create clients for
// auth, firestore, storage etc
func (f *Firebase) initializeFirebaseApp(ctx context.Context, projectAlias string, projectConfigFile string) error {
	if projectAlias == "" {
		projectAlias = defaultProject
	}
	if projectConfigFile == "" {
		projectConfigFile = firebaseProjectConfigFile
	}
	err := f.setProjectID(projectAlias, projectConfigFile)
	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}
	configs := &firebase.Config{
		ProjectID: f.projectId,
	}
	// replace this with something better
	opt := option.WithCredentialsJSON(f.credentials)
	app, err := firebase.NewApp(ctx, configs, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v\n", err)
	}
	f.App = app
	return nil
}

// Auth create a firebase auth client
func Auth(ctx context.Context, projectAlias string, projectConfigFile string) (*auth.Client, error) {
	credentials, err := constructToken()
	if err != nil {
		return nil, fmt.Errorf("Error getting credentials: %w", err)
	}
	fb := &Firebase{
		credentials: credentials,
	}
	if err := fb.initializeFirebaseApp(ctx, projectAlias, projectConfigFile); err != nil {
		return nil, err
	}
	return fb.App.Auth(ctx)
}
