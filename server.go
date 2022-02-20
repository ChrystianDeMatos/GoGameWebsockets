package main

import (
	"log"
	"net/http"
	"os"
	"websocket/handlers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer func() {
			handlers.HandleDisconnect(connection)
			connection.Close()
			log.Println("Connection closed")
		}()

		// Continuosly read and write message
		for {
			err := handlers.HandleMessages(connection)
			if err != nil {
				log.Println("read failed:", err)
				break
			}
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "game.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("is the port", port)

	http.ListenAndServe(":"+port, nil)
}
