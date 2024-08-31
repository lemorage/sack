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

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
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
				log.Printf("%sIgnoring file: %s%s", Yellow, event.Name, Reset)
				continue
			}

			if event.Op&fsnotify.Chmod != 0 {
				continue
			}

			// Use different colors for different types of file events
			switch {
			case event.Op&fsnotify.Write != 0:
				log.Printf("%sFile written: %s%s", Green, event.Name, Reset)
			case event.Op&fsnotify.Create != 0:
				log.Printf("%sFile created: %s%s", Blue, event.Name, Reset)
			case event.Op&fsnotify.Remove != 0:
				log.Printf("%sFile removed: %s%s", Red, event.Name, Reset)
			case event.Op&fsnotify.Rename != 0:
				log.Printf("%sFile renamed: %s%s", Purple, event.Name, Reset)
			}

			// Notify all connected WebSocket clients to reload
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
				if err != nil {
					log.Printf("%sWebSocket error: %v%s", Red, err, Reset)
					client.Close()
					delete(clients, client)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("%sWatcher error: %v%s", Red, err, Reset)
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
