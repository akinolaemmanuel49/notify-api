package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		URI string `yaml:"uri" envconfig:"DATABASE_URI"`
	} `yaml:"database"`
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
