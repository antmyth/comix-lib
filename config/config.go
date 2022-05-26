package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
		Port     string `yaml:"port" env:"DB_PORT" env-description:"Database port"`
		Username string `yaml:"user" env:"DB_USER" env-description:"Database user name"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-description:"Database user password"`
		Name     string `yaml:"name" env:"DB_NAME" env-description:"Database name"`
	} `yaml:"database"`
	Path   string `yaml:"path" env:"LIB_PATH" env-description:"CBZ lib path"`
	Import struct {
		ChunkSize int `yaml:"chunk" env-description:"Chunk size for cbz file processing"`
		MaxImport int `yaml:"max" env-description:"Number of issues to import per sexxion"`
	} `yaml:"import"`
}

func ReadConfig() Config {
	cfg := Config{}
	if err := cleanenv.ReadConfig("conf.yml", &cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return cfg
}
