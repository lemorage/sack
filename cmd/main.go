package main

import (
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
	mux := http.NewServeMux()
	dir := "./ui/html/pages"
	funcMap := template.FuncMap{
		"add": func(i int) int {
			return i + 1
		},
		"sub": func(i int) int {
			return i - 1
		},
	}

	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Error parsing config YAML: %s", err)
	}

	tmpl := template.Must(template.New("base").Funcs(funcMap).ParseGlob("ui/html/*.gohtml"))
	file, err := os.Open("ui/html/base.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Sort the keys of config.Pages
	var keys []string
	for key := range config.Pages {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	i := 1

	for _, page := range keys {
		pageConfig := config.Pages[page]
		pageFilename := filepath.Join(dir, fmt.Sprintf("page%d.gohtml", i))
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
			CurrentPage: i,
			TotalPages:  len(config.Pages),
			PageConfig:  pageConfig,
		})
		if err != nil {
			log.Fatalf("Error executing template for page %s: %s", page, err)
		}
		mux.HandleFunc("/model"+strconv.Itoa(i), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, pageFilename)
		})
		log.Printf("Generated HTML for %s\n", page)
		i++
	}

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)

	log.Print("Starting server on :7536")
	err = http.ListenAndServe(":7536", mux)
	log.Fatal(err)
}
