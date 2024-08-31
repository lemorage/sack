package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func parseGitignore(path string) ([]string, error) {
	var patterns []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

func isIgnored(path string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			log.Printf("Error matching pattern: %v", err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func reloadWatcher(watcher *fsnotify.Watcher, clients map[*websocket.Conn]bool, patterns []string) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if isIgnored(event.Name, patterns) {
				continue
			}

			// Log the event type and file name
			log.Printf("File event: %s, File name: %s", event.Op.String(), event.Name)

			// Check for different types of file events
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {

				// Notify all connected WebSocket clients to reload
				for client := range clients {
					err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
					if err != nil {
						log.Printf("WebSocket error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func setupWebSocket(watcher *fsnotify.Watcher) {
	clients := make(map[*websocket.Conn]bool)

	gitignorePatterns, err := parseGitignore(".gitignore")
	if err != nil {
		log.Printf("Could not parse .gitignore: %v", err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		clients[conn] = true
	})

	go reloadWatcher(watcher, clients, gitignorePatterns)
}
