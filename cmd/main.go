package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Pages map[string]PageConfig `yaml:"Pages"`
}

func main() {
	// Define command-line flags
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	port := startCmd.Int("port", 7536, "port number to start the server")

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: sack [start]")
		os.Exit(1)
	}

	// Check if the "start" command is provided
	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		if len(startCmd.Args()) > 0 {
			fmt.Println("Unexpected arguments:", startCmd.Args())
			fmt.Println("Usage: sack start [--port PORT]")
			os.Exit(1)
		}
		if startCmd.Parsed() {
			if *port < 1 || *port > 65535 {
				log.Fatalf("Invalid port number: %d. Port number must be between 1 and 65535.", *port)
				os.Exit(1)
			}

			config, err := readConfig("config.yaml")
			if err != nil {
				log.Fatalf("Error reading config file: %s", err)
			}

			tmpl := parseTemplates()
			generateHTMLFiles(config, tmpl)
			mux := setupHandlers(config)
			startServer(mux, *port)
		}
	default:
		fmt.Println("Usage: sack [start]")
		os.Exit(1)
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

func parseTemplates() *template.Template {
	funcMap := template.FuncMap{
		"add": func(i int) int {
			return i + 1
		},
		"sub": func(i int) int {
			return i - 1
		},
	}
	return template.Must(template.New("base").Funcs(funcMap).ParseGlob("ui/html/*.gohtml"))
}

func generateHTMLFiles(config Config, tmpl *template.Template) {
	dir := "./ui/html/pages"
	for i, page := range sortedPageKeys(config.Pages) {
		pageConfig := config.Pages[page]
		pageFilename := filepath.Join(dir, fmt.Sprintf("page%d.gohtml", i+1))
		newPage, err := os.Create(pageFilename)
		if err != nil {
			log.Fatalf("Error creating page file for %s: %s", page, err)
		}
		defer newPage.Close()

		err = tmpl.ExecuteTemplate(newPage, "base", struct {
			CurrentPage int
			TotalPages  int
			PageConfig  PageConfig
		}{
			CurrentPage: i + 1,
			TotalPages:  len(config.Pages),
			PageConfig:  pageConfig,
		})
		if err != nil {
			log.Fatalf("Error executing template for page %s: %s", page, err)
		}
		log.Printf("Generated HTML for %s\n", page)
	}
}

func sortedPageKeys(pages map[string]PageConfig) []string {
	keys := make([]string, 0, len(pages))
	for key := range pages {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func setupHandlers(config Config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	for i := 1; i <= len(config.Pages); i++ {
		pageFilename := fmt.Sprintf("./ui/html/pages/page%d.gohtml", i)
		mux.HandleFunc("/model"+strconv.Itoa(i), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, pageFilename)
		})
	}

	mux.HandleFunc("/", home)
	return mux
}

func startServer(mux *http.ServeMux, port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s...\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
