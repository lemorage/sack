package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"text/template"
)

// startServer initializes and runs the HTTP server on the specified port
func startServer(mux http.Handler, port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s...\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

// parseTemplates loads and parses HTML templates, adding custom template functions
func parseTemplates() *template.Template {
	funcMap := template.FuncMap{
		"add": func(i int) int { return i + 1 },
		"sub": func(i int) int { return i - 1 },
	}
	return template.Must(template.New("base").Funcs(funcMap).ParseGlob("ui/html/templates/*.gohtml"))
}

// setupHandlers configures and returns an HTTP ServeMux with all route handlers
func setupHandlers(config Config) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Serve configuration files
	mux.HandleFunc("/config.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, configPath)
	})
	mux.HandleFunc("/graph.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, storyGraphPath)
	})

	// Set up handlers for each page
	for i := 1; i <= len(config.Pages); i++ {
		pageFilename := fmt.Sprintf("./ui/html/pages/page%d.gohtml", i)
		mux.HandleFunc("/model"+strconv.Itoa(i), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, pageFilename)
		})
	}

	// Set up main route handlers
	mux.HandleFunc("/story", graph)
	mux.HandleFunc("/", home)

	return mux
}

// generateHTMLFiles creates individual HTML files for each page based on the configuration
func generateHTMLFiles(config Config, tmpl *template.Template, layout string) {
	dir := "./ui/html/pages"
	keys := sortedPageKeys(config.Pages)
	for _, key := range keys {
		pageConfig := config.Pages[key]
		pageNumber, _ := extractNumber(key)
		pageFilename := fmt.Sprintf("%s/page%d.gohtml", dir, pageNumber)

		newPage, err := os.Create(pageFilename)
		if err != nil {
			log.Fatalf("Error creating page file for %s: %s", key, err)
		}
		defer newPage.Close()

		err = tmpl.ExecuteTemplate(newPage, "base", struct {
			CurrentPage int
			TotalPages  int
			PageConfig  PageConfig
			Layout      string
		}{
			CurrentPage: pageNumber,
			TotalPages:  len(config.Pages),
			PageConfig:  pageConfig,
			Layout:      layout,
		})
		if err != nil {
			log.Fatalf("Error executing template for page %s: %s", key, err)
		}
		log.Printf("Generated HTML for %s\n", key)
	}
}

// sortedPageKeys returns a sorted list of page keys from the configuration
func sortedPageKeys(pages map[string]PageConfig) []string {
	keys := make([]int, 0, len(pages))
	keyMap := make(map[int]string)
	for key := range pages {
		pageNumber, err := extractNumber(key)
		if err != nil {
			log.Fatalf("Error: Key '%s' does not contain a number", key)
		}
		keys = append(keys, pageNumber)
		keyMap[pageNumber] = key
	}
	sort.Ints(keys)
	sortedKeys := make([]string, len(keys))
	for i, num := range keys {
		sortedKeys[i] = keyMap[num]
	}
	return sortedKeys
}

// extractNumber retrieves the numeric part from a string key
func extractNumber(key string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(key)
	if numStr == "" {
		return 0, fmt.Errorf("no number found in key")
	}
	return strconv.Atoi(numStr)
}
