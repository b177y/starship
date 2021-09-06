package main

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Type   string `yaml:"type"`
		Source string `yaml:"src"`
	} `yaml:"db"`
	Quasar struct {
		Name   string `yaml:"name"`
		Listen struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"listen"`
	} `yaml:"quasar"`
	Authsecret string
}

// https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp
func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
