package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const description = `Kamanda is a  Firebase CLI Tool extender and should be used alongside it. 

Kamanda provides additional functionality currently not available to via the 
Firebase CLI Tool such as User Management, Cloud Firestore Management etc from the CLI.

For instance, it allows you to easily create users with custom tokens, 
which is always a trick preposition.`

var cfgFile string

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kamanda.yaml)")
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
		viper.AddConfigPath(home)
		viper.SetConfigFile(".kamanda/refresh_token.json")
	}
	viper.Set("GOOGLE_OAUTH_CLIENT_ID", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	viper.Set("GOOGLE_OAUTH_CLIENT_ID", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
