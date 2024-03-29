// Package config provides exports for reading in configuration files.
package config

import (
	"io"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

// Config is the data structure that represents the root of the configuration.
type Config struct {
	ColorVariant string `yaml:"colorVariant"`
	Modules      []any  `yaml:"modules"`
}

var defaultConfig = Config{
	ColorVariant: "dark",
	Modules: []any{
		map[any]any{
			"module":  "network",
			"pattern": "(en|wl)+",
		},
		map[any]any{
			"module":    "datetime",
			"format":    time.RFC1123,
			"timezones": []string{"Local"},
			"interval":  1,
		},
		map[any]any{
			"module":  "text",
			"content": "gobar",
		},
	},
}

// GetConfig will retrieve a configuration file at an optionally overridden
// location.
func GetConfig(overrideFilepath string) (*Config, error) {
	config := defaultConfig

	path, err := getConfigFilePath(overrideFilepath)
	if err == ErrNoLookupLocation || err == ErrNoConfig {
		return &config, nil
	}
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
