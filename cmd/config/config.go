package config

import (
	"bytes"
	_ "embed"
	"github.com/spf13/viper"
	"strings"
)

//go:embed config.yaml
var defaultConfig []byte

type Postgres struct {
	Host     string
	Port     uint
	Username string
	Password string
}

type Config struct {
	Postgres Postgres
}

func Read() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
