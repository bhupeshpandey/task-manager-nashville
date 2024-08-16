package config

import (
	"github.com/bhupeshpandey/task-manager-nashville/internal/models"
	yaml "gopkg.in/yaml.v2"
	"os"
)

// LoadConfig reads the YAML configuration file and unmarshals it into a Config struct
func LoadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
