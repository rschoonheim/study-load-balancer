package configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"hash/crc32"
	"log/slog"
	"os"
	"time"
)

// checksum - checksum of the configuration file
var previousChecksum uint32

// calculateChecksum - calculates the checksum of the configuration file
func calculateChecksum(configuration *Configuration) uint32 {
	// json encoding of the configuration
	//
	str, err := json.Marshal(configuration)
	if err != nil {
		slog.Error("Failed to marshal configuration.", "error", err.Error())
		return 0
	}

	// Quick hashing of the configuration
	//
	data := crc32.Checksum([]byte(str), crc32.IEEETable)

	// Calculate the checksum of the configuration
	//
	return data
}

// loadConfigurationFromFilesystem - loads the configuration from the filesystem
func loadConfigurationFromFilesystem(configurationPath string) *Configuration {

	// Check if the main configuration file exists
	//
	_, err := os.Stat(configurationPath)
	if os.IsNotExist(err) {
		println("Main configuration file does not exist")
		os.Exit(1)
	}

	// Read contents of the main configuration file
	//
	mainConfigurationFile, err := os.Open(configurationPath)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// Parse the main configuration file
	//
	var configuration *Configuration
	if err := yaml.NewDecoder(mainConfigurationFile).Decode(&configuration); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	return configuration
}

// Watcher - configuration watcher scans the configuration file for changes and updates the configuration
// accordingly.
func Watcher(configuration *Configuration, configurationPath string) {
	// Parse the wait time for the configuration watcher
	//
	rescan, err := time.ParseDuration(configuration.Filesystem.Configurations.Rescan)
	if err != nil {
		slog.Error("Failed to parse rescan duration.", "error", err.Error())
		return
	}
	slog.Info("Scanning configurations for changes.")

	// Load the current configuration from the filesystem
	//
	newConfiguration := loadConfigurationFromFilesystem(configurationPath)

	// Calculate the checksum of the configuration
	//
	newChecksum := calculateChecksum(newConfiguration)
	currentChecksum := calculateChecksum(configuration)
	if newChecksum == currentChecksum {
		slog.Info("No changes detected in the configuration file.", "current_checksum", newChecksum, "previous_checksum", currentChecksum)
		time.Sleep(rescan)
		Watcher(configuration, configurationPath)
	}

	previousChecksum = newChecksum

	// Update the configuration with the new configuration
	//
	*configuration = *newConfiguration
	slog.Info("Configuration updated.", "current_checksum", newChecksum, "previous_checksum", currentChecksum)

	Watcher(configuration, configurationPath)
}
