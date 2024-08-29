package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var configPath = "configs/config.yaml"
var storyGraphPath = "configs/graph.json"
var pathsToWatch = []string{
	"./ui/html",
	"./ui/static",
	"./configs",
	"./cmd/main.go",
	"./cmd/handlers.go",
	"./cmd/middleware.go",
}

func main() {
	// Initialize file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, path := range pathsToWatch {
		err := addPathsRecursively(watcher, path)
		if err != nil {
			log.Fatalf("Error setting up watcher for path %s: %v", path, err)
		}
	}

	// Set up WebSocket server for auto-reload
	setupWebSocket(watcher)

	// Define command-line flags
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	port := startCmd.Int("port", 7536, "port number to start the server")
	layout := startCmd.String("layout", "card", "layout of the pages (card or plain)")

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
			fmt.Println("Usage: sack start [--port PORT] [--layout LAYOUT]")
			os.Exit(1)
		}
		if startCmd.Parsed() {
			// Validate port number
			if *port < 1 || *port > 65535 {
				log.Fatalf("Invalid port number: %d. Port number must be between 1 and 65535.", *port)
				os.Exit(1)
			}

			// Validate layout
			if *layout != "card" && *layout != "plain" {
				log.Fatalf("Invalid layout: %s. Layout must be either 'card' or 'plain'.", *layout)
				os.Exit(1)
			}

			config, err := readConfig(configPath)
			if err != nil {
				log.Fatalf("Error reading config file: %s", err)
			}

			tmpl := parseTemplates()
			generateHTMLFiles(config, tmpl, *layout)
			mux := setupHandlers(config)

			// Add the middleware to inject the WebSocket script
			muxWithMiddleware := injectWebSocketScriptMiddleware(mux)

			startServer(muxWithMiddleware, *port)
		}
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if generateCmd.Parsed() {
			if *batch > 0 {
				if *batch > 1024 {
					log.Fatalf("Bulk number too large: %d. Must be between 1 and 1024.", *batch)
					os.Exit(1)
				}
				config, err := readConfig(configPath)
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

// interactiveGenerate prompts the user for input to create a new page configuration
func interactiveGenerate() {
	reader := bufio.NewReader(os.Stdin)
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Prompt user for various page configuration details
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

	// Create new PageConfig with user input
	pageConfig := PageConfig{
		ModelSrcPath:    modelSrcPath,
		ModelIosSrcPath: modelIosSrcPath,
		PosterPath:      posterPath,
		Description:     description,
		ModelName:       modelName,
		DesignerWebsite: designerWebsite,
		DesignerName:    designerName,
	}

	// Add new page to config
	pageName := fmt.Sprintf("page%d", len(config.Pages)+1)
	config.Pages[pageName] = pageConfig

	// Save updated config
	writeConfig(configPath, config)
}

// batchGenerate creates multiple new page configurations automatically
func batchGenerate(config Config, count int) {
	pageCount := len(config.Pages)

	// Generate 'count' number of new pages
	for i := 1; i <= count; i++ {
		pageConfig := PageConfig{
			ModelSrcPath:    fmt.Sprintf("/static/obj%d/object%d.glb", pageCount+i, pageCount+i),
			ModelIosSrcPath: fmt.Sprintf("/static/obj%d/object%d.usdz", pageCount+i, pageCount+i),
			PosterPath:      fmt.Sprintf("/static/obj%d/object%d.webp", pageCount+i, pageCount+i),
			Description:     fmt.Sprintf("This is my masterpiece %d", pageCount+i),
			ModelName:       fmt.Sprintf("Model %d", pageCount+i),
			DesignerWebsite: config.Pages["page1"].DesignerWebsite, // Use existing values for designer info
			DesignerName:    config.Pages["page1"].DesignerName,
		}

		// Add new page to config
		pageName := fmt.Sprintf("page%d", pageCount+i)
		config.Pages[pageName] = pageConfig
	}

	writeConfig(configPath, config)
}
