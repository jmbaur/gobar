package config

import (
	"io"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Modules []any `yaml:"modules"`
}

func GetConfig(flagConfigFile string) (*Config, error) {
	config := Config{
		Modules: []any{
			map[any]any{
				"module":  "network",
				"pattern": "(en|wl)+",
			},
			map[any]any{
				"module":   "datetime",
				"format":   time.RFC1123,
				"interval": 1,
			},
			map[any]any{
				"module":  "text",
				"content": "gobar",
			},
		},
	}

	path, err := getConfigFilePath(flagConfigFile)
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
