package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"websocket/handlers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var todoList []string

func getCmd(input string) string {
	inputArr := strings.Split(input, " ")
	return inputArr[0]
}

func getMessage(input string) string {
	inputArr := strings.Split(input, " ")
	var result string
	for i := 1; i < len(inputArr); i++ {
		result += inputArr[i]
	}
	return result
}

func updateTodoList(input string) {
	tmpList := todoList
	todoList = []string{}
	for _, val := range tmpList {
		if val == input {
			continue
		}
		todoList = append(todoList, val)
	}
}

var playersConnections []*websocket.Conn

func notifyAllConnections(connections []*websocket.Conn) {
	for _, conn := range connections {
		conn.WriteMessage(websocket.TextMessage, []byte("New Todo Added"))
	}
}

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
