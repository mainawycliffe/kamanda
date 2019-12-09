package firebase

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	s "github.com/bitfield/script"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"log"
)

const firebaseProjectConfig string = "./.firebaserc"
const defaultProject = "default"

type Firebase struct {
	App       *firebase.App
	projectId string
}

func (f *Firebase) setProjectID(projectId string) error {

	configFileContent, err := s.File(firebaseProjectConfig).Bytes()

	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}

	var decodedConfigs *FirebaseProjectConfigs

	err = json.Unmarshal(configFileContent, &decodedConfigs)

	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}

	f.projectId = decodedConfigs.Projects[projectId]

	return nil
}

func (f *Firebase) InitializeFirbeaseApp(ctx context.Context, projectId string) error {

	if projectId == "" {
		projectId = defaultProject
	}

	err := f.setProjectID(projectId)

	if err != nil {
		return fmt.Errorf("An error occurred while reading config file: %w", err)
	}

	configs := &firebase.Config{
		ProjectID: f.projectId,
	}

	// replace this with something better
	opt := option.WithCredentialsFile(viper.GetString("refreshTokenFilePath"))

	app, err := firebase.NewApp(ctx, configs, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	f.App = app

	return nil
}

func (f *Firebase) Auth(ctx context.Context) (*auth.Client, error) {
	return f.App.Auth(ctx)
}
