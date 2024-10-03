package main

import (
	"gopkg.in/yaml.v2"
	configuration2 "load-balancer/configuration"
	"os"
)

var configuration configuration2.Configuration

func init() {

	// Check if `MAIN_CONFIGURATION` is set
	//
	mainConfigurationPath := os.Getenv("MAIN_CONFIGURATION")
	if mainConfigurationPath == "" {
		println("MAIN_CONFIGURATION environment variable is not set")
		os.Exit(1)
	}

	// Check if the main configuration file exists
	//
	_, err := os.Stat(mainConfigurationPath)
	if os.IsNotExist(err) {
		println("Main configuration file does not exist")
		os.Exit(1)
	}

	// Read contents of the main configuration file
	//
	mainConfigurationFile, err := os.Open(mainConfigurationPath)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// Parse the main configuration file
	//
	if err := yaml.NewDecoder(mainConfigurationFile).Decode(&configuration); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func main() {
	println(configuration.Filesystem.Configurations.Rescan)
}
