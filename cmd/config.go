package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// PageConfig is a struct that holds the configuration for a page
type PageConfig struct {
	ModelSrcPath    string `yaml:"ModelSrcPath"`
	ModelIosSrcPath string `yaml:"ModelIosSrcPath"`
	PosterPath      string `yaml:"PosterPath"`
	Description     string `yaml:"Description"`
	ModelName       string `yaml:"ModelName"`
	DesignerWebsite string `yaml:"DesignerWebsite"`
	DesignerName    string `yaml:"DesignerName"`
}

type Config struct {
	Pages map[string]PageConfig `yaml:"Pages"`
}

func writeConfig(filename string, config Config) {
	configData, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Error marshalling config: %s", err)
	}

	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err = file.Write(configData); err != nil {
		log.Fatal(err)
	}
}

func readConfig(filename string) (Config, error) {
	var config Config
	configData, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(configData, &config)
	return config, err
}
