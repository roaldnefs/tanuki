// Copyright © 2018 Roald Nefs <info@roaldnefs.com>

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

const (
	defaultBaseURL = "https://gitlab.com/"
)

var (
	// Used for flags.
	cfgFile, token, baseURL   string
	enableDryRun, enableDebug bool

	// git represent the GitLab API client.
	git *gitlab.Client

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "tanuki",
		Short: "A tool for performing actions on GitLab repos or a single repo.",
		Long:  `A tool for performing actions on GitLab repos or a single repo.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "d", false, "enable debug logging (default false)")
	rootCmd.PersistentFlags().BoolVarP(&enableDryRun, "dry-run", "", false, "do not change settings just print the changes that would occur (default false)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tanuki.yaml)")

	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitLab API token")
	rootCmd.MarkFlagRequired("token")
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	// Base URL for APU requests. Defaults to the public GitLab API, but can be
	// set to a domain endpoint to use with a self hosted GitLab Server. baseURL
	// should always be specified with a trailing slash.
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", defaultBaseURL, "GitLab URL")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))

	// Initialize the GitLab client.
	git = gitlab.NewClient(nil, token)
	git.SetBaseURL(baseURL)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tanuki" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tanuki")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
