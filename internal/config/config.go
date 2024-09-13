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
	MigrationsDir string `yaml:"migrations-dir" env:"MIGRATIONS_DIR" env-default:"./migrations/prod"`
	Conn          string `yaml:"conn" env:"POSTGRES_CONN" env-default:""`
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
