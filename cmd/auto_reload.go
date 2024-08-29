package main

import (
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reloadWatcher(watcher *fsnotify.Watcher, clients map[*websocket.Conn]bool) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Log the event type and file name
			log.Printf("File event: %s, File name: %s", event.Op.String(), event.Name)

			// Check for different types of file events
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {

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

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		clients[conn] = true
	})

	go reloadWatcher(watcher, clients)
}
