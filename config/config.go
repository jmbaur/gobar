package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Modules []interface{} `yaml:"modules"`
}

var ErrNoLookupLocation = errors.New("No config file lookup location")

func getConfigFilePath(flagConfigFile string) (string, error) {
	if flagConfigFile != "" {
		return filepath.Join(flagConfigFile), nil
	}

	xdgConfigHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if ok {
		return filepath.Join(xdgConfigHome, "gobar", "gobar.yaml"), nil
	}

	home, ok := os.LookupEnv("HOME")
	if ok {
		return filepath.Join(home, ".config", "gobar", "gobar.yaml"), nil
	}

	return "", ErrNoLookupLocation
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
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return &config, nil
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
