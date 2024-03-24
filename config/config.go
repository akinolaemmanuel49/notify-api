package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		URI string `yaml:"uri"`
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
