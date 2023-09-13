package config

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Acme  Acme  `yaml:"acme"`
	Proxy Proxy `yaml:"proxy"`
}

var instance Config
var once = sync.Once{}

func Get() Config {
	once.Do(func() {
		instance = loadConfig()
		log.Printf("%#v", instance)
	})

	return instance
}

func loadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Println(err)
		}
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yml"
	}

	config := DefaultConfig

	if err := fromYaml(path, &config); err != nil {
		log.Printf("couldn'n load config from %s: %s\r\n", path, err.Error())
	}

	if err := fromEnv(&config); err != nil {
		log.Printf("couldn'n load config from env: %s\r\n", err.Error())
	}

	return config
}

func fromYaml(path string, config *Config) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func fromEnv(config *Config) error {
	return envconfig.Process("", config)
}
