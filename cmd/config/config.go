package config

import (
	"bytes"
	_ "embed"
	"github.com/spf13/viper"
	"log"
	"strings"
)

//go:embed config.yaml
var defaultConfig []byte

type Database struct {
	Host     string
	Port     uint
	Username string
	Password string
	Dbname   string
}

type Server struct {
	Port string
}
type Config struct {
	Database Database
	Server   Server
}

func Read() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	host := viper.GetString("database.host")
	log.Printf("host from env: %s", host)

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		return nil, err
	}

	host = viper.GetString("database.host")
	log.Printf("host from config: %s", host)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
