/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/si-bas/go-rest-geospatial/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	ConsulURL   = "CONSUL_URL"
	ConsulToken = "CONSUL_TOKEN"
	ConsulKey   = "CONSUL_KEY"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geospatial-service",
	Short: "Geospatial Service API",
	Long:  `Geospatial Service API to serve geospatial data.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is in current directory .env)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in config/files directory with name ".config.json".
		viper.SetConfigType("json")
		viper.AddConfigPath("./")
		viper.SetConfigName(".config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		panic("error loading configuration: no config")
	}

	if err := viper.Unmarshal(&config.Config); err != nil {
		fmt.Fprintln(os.Stderr, "failed to unmarshal config to struct variable", err.Error())
	}
}
