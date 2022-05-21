package config

import (
	"io"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Modules []interface{} `yaml:"modules"`
}

func GetConfig(flagConfigFile string) (*Config, error) {
	config := Config{
		Modules: []interface{}{
			map[interface{}]interface{}{
				"module": "datetime",
				"format": time.RFC1123,
			},
		},
	}

	path, err := getConfigFilePath(flagConfigFile)
	if err == ErrNoLookupLocation {
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
