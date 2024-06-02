package main

import (
	"os"
	"reflect"
	"testing"
)

func TestWriteConfig(t *testing.T) {
	// Prepare test data
	config := Config{
		Pages: map[string]PageConfig{
			"page1": {
				ModelSrcPath:    "model1.glb",
				ModelIosSrcPath: "model1.usdz",
				PosterPath:      "poster1.png",
				Description:     "Test Model 1",
				ModelName:       "Model 1",
				DesignerWebsite: "https://example.com/designer1",
				DesignerName:    "Designer 1",
			},
		},
	}

	// Write config to a temporary file
	filename := "test_config.yaml"
	defer os.Remove(filename)
	writeConfig(filename, config)

	// Read the file back and compare
	readConfig, err := readConfig(filename)
	if err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	if !reflect.DeepEqual(config, readConfig) {
		t.Fatalf("Expected %v, got %v", config, readConfig)
	}
}

func TestReadConfig(t *testing.T) {
	// Prepare a test config file
	configData := `
Pages:
  page1:
    ModelSrcPath: "model1.glb"
    ModelIosSrcPath: "model1.usdz"
    PosterPath: "poster1.png"
    Description: "Test Model 1"
    ModelName: "Model 1"
    DesignerWebsite: "https://example.com/designer1"
    DesignerName: "Designer 1"
`
	filename := "test_config.yaml"
	defer os.Remove(filename)
	os.WriteFile(filename, []byte(configData), 0644)

	// Read the config
	config, err := readConfig(filename)
	if err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	expected := Config{
		Pages: map[string]PageConfig{
			"page1": {
				ModelSrcPath:    "model1.glb",
				ModelIosSrcPath: "model1.usdz",
				PosterPath:      "poster1.png",
				Description:     "Test Model 1",
				ModelName:       "Model 1",
				DesignerWebsite: "https://example.com/designer1",
				DesignerName:    "Designer 1",
			},
		},
	}

	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("Expected %v, got %v", expected, config)
	}
}
