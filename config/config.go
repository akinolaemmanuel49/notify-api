package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		USER string `yaml:"user" envconfig:"DB_USER"`
		PASS string `yaml:"pass" envconfig:"DB_PASS"`
		NAME string `yaml:"name" envconfig:"DB_NAME"`
	} `yaml:"database"`
	JWT struct {
		KEY       string `yaml:"key" envconfig:"JWT_KEY"`
		TOKEN_TTL string `yaml:"tokenTTL" envconfig:"JWT_TOKEN_TTL"`
	} `yaml:"jwt"`
}

func processError(err error) {
	log.Fatalln(err)
}
func (cfg *Config) ReadFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func (cfg *Config) ReadEnv() {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
