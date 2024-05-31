package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	batch := generateCmd.Int("batch", 0, "generate multiple pages in batch")

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: sack [start | generate]")
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
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if generateCmd.Parsed() {
			if *batch > 0 {
				if *batch > 1024 {
					log.Fatalf("Bulk number too large: %d. Must be between 1 and 1024.", *batch)
					os.Exit(1)
				}
				config, err := readConfig("config.yaml")
				if err != nil || len(config.Pages) == 0 {
					log.Fatalf("Error reading config file or no existing pages to reference: %s", err)
					os.Exit(1)
				}
				batchGenerate(config, *batch)
			} else if len(os.Args[2:]) == 0 {
				interactiveGenerate()
			} else {
				fmt.Println("Usage: sack generate [--batch num]")
				os.Exit(1)
			}
		}
	default:
		fmt.Println("Usage: sack [start | generate]")
		os.Exit(1)
	}
}

func interactiveGenerate() {
	reader := bufio.NewReader(os.Stdin)
	config, err := readConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	fmt.Print("Enter ModelSrcPath: ")
	modelSrcPath, _ := reader.ReadString('\n')
	modelSrcPath = strings.TrimSpace(modelSrcPath)

	fmt.Print("Enter ModelIosSrcPath: ")
	modelIosSrcPath, _ := reader.ReadString('\n')
	modelIosSrcPath = strings.TrimSpace(modelIosSrcPath)

	fmt.Print("Enter PosterPath: ")
	posterPath, _ := reader.ReadString('\n')
	posterPath = strings.TrimSpace(posterPath)

	fmt.Print("Enter Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Enter ModelName: ")
	modelName, _ := reader.ReadString('\n')
	modelName = strings.TrimSpace(modelName)

	fmt.Print("Enter DesignerWebsite: ")
	designerWebsite, _ := reader.ReadString('\n')
	designerWebsite = strings.TrimSpace(designerWebsite)

	fmt.Print("Enter DesignerName: ")
	designerName, _ := reader.ReadString('\n')
	designerName = strings.TrimSpace(designerName)

	pageConfig := PageConfig{
		ModelSrcPath:    modelSrcPath,
		ModelIosSrcPath: modelIosSrcPath,
		PosterPath:      posterPath,
		Description:     description,
		ModelName:       modelName,
		DesignerWebsite: designerWebsite,
		DesignerName:    designerName,
	}

	pageName := fmt.Sprintf("page%d", len(config.Pages)+1)
	config.Pages[pageName] = pageConfig

	writeConfig("config.yaml", config)
}

func batchGenerate(config Config, count int) {
	pageCount := len(config.Pages)
	for i := 1; i <= count; i++ {
		pageConfig := PageConfig{
			ModelSrcPath:    fmt.Sprintf("/static/obj%d/object%d.glb", pageCount+i, pageCount+i),
			ModelIosSrcPath: fmt.Sprintf("/static/obj%d/object%d.usdz", pageCount+i, pageCount+i),
			PosterPath:      fmt.Sprintf("/static/obj%d/object%d.webp", pageCount+i, pageCount+i),
			Description:     fmt.Sprintf("This is my masterpiece %d", pageCount+i),
			ModelName:       fmt.Sprintf("Model %d", pageCount+i),
			DesignerWebsite: config.Pages["page1"].DesignerWebsite,
			DesignerName:    config.Pages["page1"].DesignerName,
		}
		pageName := fmt.Sprintf("page%d", pageCount+i)
		config.Pages[pageName] = pageConfig
	}

	writeConfig("config.yaml", config)
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
	keys := sortedPageKeys(config.Pages)

	for _, key := range keys {
		pageConfig := config.Pages[key]
		pageNumber, _ := extractNumber(key)
		pageFilename := filepath.Join(dir, fmt.Sprintf("page%d.gohtml", pageNumber))
		newPage, err := os.Create(pageFilename)
		if err != nil {
			log.Fatalf("Error creating page file for %s: %s", key, err)
		}
		defer newPage.Close()

		err = tmpl.ExecuteTemplate(newPage, "base", struct {
			CurrentPage int
			TotalPages  int
			PageConfig  PageConfig
		}{
			CurrentPage: pageNumber,
			TotalPages:  len(config.Pages),
			PageConfig:  pageConfig,
		})
		if err != nil {
			log.Fatalf("Error executing template for page %s: %s", key, err)
		}
		log.Printf("Generated HTML for %s\n", key)
	}
}

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

func extractNumber(key string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(key)
	if numStr == "" {
		return 0, fmt.Errorf("no number found in key")
	}
	return strconv.Atoi(numStr)
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
