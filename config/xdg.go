package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrNoLookupLocation represents when there is no filepath location to
	// find a configuration file.
	ErrNoLookupLocation = errors.New("no config file lookup location")
	// ErrNoConfig represents when no configuration file was loaded.
	ErrNoConfig = errors.New("no config file loaded")
)

func configFilePriority() []string {
	dirs := []string{}

	if xdgConfigHome, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		dirs = append(dirs, filepath.Join(xdgConfigHome, "gobar", "gobar.yaml"))
	}
	if home, ok := os.LookupEnv("HOME"); ok {
		dirs = append(dirs, filepath.Join(home, ".config", "gobar", "gobar.yaml"))
	}

	if xdgConfigDirs, ok := os.LookupEnv("XDG_CONFIG_DIRS"); ok {
		for _, dir := range strings.Split(xdgConfigDirs, ":") {
			dirs = append(dirs, filepath.Join(dir, "gobar", "gobar.yaml"))
		}
	}

	dirs = append(dirs, filepath.Join("etc", "xdg", "gobar", "gobar.yaml"))

	return dirs
}

func getConfigFilePath(flagConfigFile string) (string, error) {
	if flagConfigFile == "NONE" {
		return "", ErrNoConfig
	}

	if flagConfigFile != "" {
		return filepath.Abs(filepath.Join(flagConfigFile))
	}

	for _, location := range configFilePriority() {
		if _, err := os.Stat(location); err == nil {
			return location, nil
		}
	}

	return "", ErrNoLookupLocation
}
