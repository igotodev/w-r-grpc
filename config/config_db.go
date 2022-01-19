package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"my_db"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
}

func InitConfig() PostgresConfig {
	var cfg PostgresConfig

	if err := cleanenv.ReadConfig("config_db.yaml", &cfg); err != nil {
		log.Println(err)
	}

	return cfg
}
