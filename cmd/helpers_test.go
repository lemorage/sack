package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

func TestStartServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func TestGenerateHTMLFiles(t *testing.T) {
	// Create a temporary directory for HTML files
	dir := "./ui/html/pages"
	os.MkdirAll(dir, 0755)
	defer os.RemoveAll(dir)

	// Create test config and template
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

	tmpl := template.Must(template.New("base").Parse("Page: {{.PageConfig.ModelName}}"))

	// Generate HTML files
	generateHTMLFiles(config, tmpl)

	// Check if the file is created
	pageFilename := filepath.Join(dir, "page1.gohtml")
	if _, err := os.Stat(pageFilename); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to be created", pageFilename)
	}
}

func TestSetupHandlers(t *testing.T) {
	// Setup the temporary directory and file structure
	err := os.MkdirAll("./ui/html/pages", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}
	defer os.RemoveAll("./ui") // Clean up after the test

	// Create a dummy page1.gohtml file
	pageFilePath := "./ui/html/pages/page1.gohtml"
	pageFileContent := "<html><body>Test Page</body></html>"
	err = os.WriteFile(pageFilePath, []byte(pageFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create page file: %v", err)
	}

	// Create a config with a single page
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

	mux := setupHandlers(config)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/model1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}
