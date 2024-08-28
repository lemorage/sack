package main

import (
	"fmt"
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
			// If a file changes, send reload signal
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("File modified:", event.Name)
				for client := range clients {
					err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
					if err != nil {
						log.Println("Write error:", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
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
