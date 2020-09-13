package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type ServerConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func Parse(configPath string) (cfg *Config, err error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	return
}
