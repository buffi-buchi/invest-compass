package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfig(path string) (Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return Config{}, fmt.Errorf("open config: %w", err)
	}

	defer file.Close()

	var config Config

	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	if err = config.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return config, nil
}

type Config struct {
	Server      ServerConfig `yaml:"server"`
	DebugServer ServerConfig `yaml:"debug_server"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("validate server config: %w", err)
	}

	if err := c.DebugServer.Validate(); err != nil {
		return fmt.Errorf("validate debug server config: %w", err)
	}

	return nil
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

func (c ServerConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}

	return nil
}
