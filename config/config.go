package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		User string `yaml:"user" envconfig:"DB_USER"`
		Pass string `yaml:"pass" envconfig:"DB_PASS"`
		Name string `yaml:"name" envconfig:"DB_NAME"`
	} `yaml:"database"`
	JWT struct {
		Key       string `yaml:"key" envconfig:"JWT_KEY"`
		Token_TTL string `yaml:"tokenTTL" envconfig:"JWT_TOKEN_TTL"`
	} `yaml:"jwt"`
	RateLimiting struct {
		MaxRequests string `yaml:"maxRequests" envconfig:"MAX_REQUESTS"`
		Duration    string `yaml:"duration" envconfig:"REQUEST_LIMIT_DURATION"`
	} `yaml:"rateLimiting"`
}

func processError(err error) {
	log.Println(err)
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
