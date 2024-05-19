package main

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
