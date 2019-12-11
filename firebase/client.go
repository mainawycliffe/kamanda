package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	s "github.com/bitfield/script"
	"google.golang.org/api/option"
)

const (
	firebaseProjectConfig = "./.firebaserc"
	defaultProject        = "default"
)

type Firebase struct {
	App       *firebase.App
	projectId string
}

// setProjectID use the project alias to get the firebase project id
func (f *Firebase) setProjectID(projectAlias string) error {
	configFileContent, err := s.File(firebaseProjectConfig).Bytes()
	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}
	var decodedConfigs *FirebaseProjectConfigs
	err = json.Unmarshal(configFileContent, &decodedConfigs)
	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}
	f.projectId = decodedConfigs.Projects[projectAlias]
	return nil
}

// initializeFirbeaseApp create a new firebase app that can create clients for
// auth, firestore, storage etc
func (f *Firebase) initializeFirbeaseApp(ctx context.Context, projectAlias string) error {
	if projectAlias == "" {
		projectAlias = defaultProject
	}
	err := f.setProjectID(projectAlias)
	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}
	configs := &firebase.Config{
		ProjectID: f.projectId,
	}
	credentials, err := constructToken()
	if err != nil {
		return fmt.Errorf("Error getting credentials: %w", err)
	}
	// replace this with something better
	opt := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(ctx, configs, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	f.App = app
	return nil
}

// Auth create a firebase auth client
func Auth(ctx context.Context, projectAlias string) (*auth.Client, error) {
	fb := &Firebase{}
	if err := fb.initializeFirbeaseApp(ctx, projectAlias); err != nil {
		return nil, err
	}
	return fb.App.Auth(ctx)
}
