package cmd

import (
	"fmt"
	"github.com/mainawycliffe/kamanda/configs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// "In this context, the client secret is obviously not treated as a secret"
// https://developers.google.com/identity/protocols/OAuth2InstalledApp
const (
	GOOGLE_OAUTH_CLIENT_ID     = "69911959250-cbimd6dvmvgqlt45tgqf8kkmo5br5vql.apps.googleusercontent.com"
	GOOGLE_OAUTH_CLIENT_SECRET = "-U9Ab0SHVO1MuuES2CkHAB1e"
)

const description = `Kamanda is a  Firebase CLI Tool extender and should be used alongside it. 

Kamanda provides additional functionality currently not available to via the 
Firebase CLI Tool such as User Management, Cloud Firestore Management etc from the CLI.

For instance, it allows you to easily create users with custom tokens, 
which is always a trick preposition.`

var (
	cfgFile         string
	firebaseToken   string
	firebaseProject string
)

var rootCmd = &cobra.Command{
	Use:   "kamanda",
	Short: "Kamanda is an extender Firebase Tools CLI",
	Long:  description,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kamanda/config.json)")
	rootCmd.PersistentFlags().StringVar(&firebaseToken, "token", "", "firebase token to use for authentication")
	// this can be used to pass project alias to sub commands, incase having
	// multiple projects
	rootCmd.PersistentFlags().StringVarP(&firebaseProject, "project", "P", "default", "The firebase project to use")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configPath := fmt.Sprintf("%s/.kamanda.yaml", home)
		_, err = os.Stat(configPath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error checking if config file exists: %s\n", err.Error())
			os.Exit(1)
		}
		viper.SetConfigFile(configPath)
	}
	viper.Set(configs.GoogleOAuthClientIDConfigKey, GOOGLE_OAUTH_CLIENT_ID)
	viper.Set(configs.GoogleOAuthClientSecretConfigKey, GOOGLE_OAUTH_CLIENT_SECRET)
	viper.AutomaticEnv()
	// @todo: improve error handling here, i.e fail if error is due to missing
	// config file
	_ = viper.SafeWriteConfig()
	// @todo: probably bind the token flag to the refresh setting config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading configs: %s\n", err.Error())
		os.Exit(1)
	}
}
