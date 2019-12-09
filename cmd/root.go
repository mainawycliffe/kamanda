package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const configDirName = ".kamanda"

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
		configPath := fmt.Sprintf("%s/%s/configs.json", home, configDirName)
		configDirPath := fmt.Sprintf("%s/%s", home, configDirName)
		_, err = os.Stat(configPath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error checking if config file exists: %s\n", err.Error())
			os.Exit(1)
		}
		if os.IsNotExist(err) {
			err := os.MkdirAll(configDirPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error checking if config file exists: %s\n", err.Error())
				os.Exit(1)
			}
		}
		viper.SetConfigFile(configPath)
	}
	viper.Set("GOOGLE_OAUTH_CLIENT_ID", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	viper.Set("GOOGLE_OAUTH_CLIENT_SECRET", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	viper.AutomaticEnv()
	err := viper.SafeWriteConfig()
	if err != nil {
		fmt.Printf("Error checking if config file exists: %s\n", err.Error())
		os.Exit(1)
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
