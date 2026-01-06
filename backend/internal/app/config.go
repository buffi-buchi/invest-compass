package app

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func ReadConfig(path string) (Config, error) {
	config, err := readConfig(path)
	if err != nil {
		return Config{}, err
	}

	if err = config.validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return config, nil
}

func readConfig(path string) (Config, error) {
	v := viper.New()

	v.SetConfigFile(path)

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var config Config

	if err := v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}

type Config struct {
	Server      ServerConfig   `mapstructure:"server"`
	DebugServer ServerConfig   `mapstructure:"debugServer"`
	Postgres    PostgresConfig `mapstructure:"postgres"`
}

func (c Config) validate() error {
	if err := c.Server.validate(); err != nil {
		return fmt.Errorf("validate server config: %w", err)
	}

	if err := c.DebugServer.validate(); err != nil {
		return fmt.Errorf("validate debug server config: %w", err)
	}

	if err := c.Postgres.validate(); err != nil {
		return fmt.Errorf("validate postgres config: %w", err)
	}

	return nil
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func (c ServerConfig) validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}

	return nil
}

type PostgresConfig struct {
	Database string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

func (c PostgresConfig) validate() error {
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if c.User == "" {
		return fmt.Errorf("database user is required")
	}

	if c.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Port == "" {
		return fmt.Errorf("database port is required")
	}

	return nil
}

func (c PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
