package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type ServerConfig struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:":8080"`
}

type PostgresConfig struct {
	Conn     string `yaml:"conn" env:"POSTGRES_CONN" env-default:""`
	JdbcUrl  string `yaml:"jdbc-url" env:"POSTGRES_JDBC_URL" env-default:""`
	Username string `yaml:"username" env:"POSTGRES_USERNAME" env-default:""`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:""`
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:""`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:""`
	Database string `yaml:"database" env:"POSTGRES_DATABASE" env-default:""`
}

func MustLoad(configPath string) Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Cannot find config file")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Error while reading config")
	}

	return cfg
}
