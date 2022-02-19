package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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

func notifyAllConnections(connections []*websocket.Conn) {
	for _, conn := range connections {
		conn.WriteMessage(websocket.TextMessage, []byte("New Todo Added"))
	}
}

func main() {

	connections := []*websocket.Conn{}

	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		connections = append(connections, connection)
		defer connection.Close()

		// Continuosly read and write message
		for {
			messageType, message, err := connection.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}
			userInput := string(message)
			command := getCmd(userInput)
			messageOfTodo := getMessage(userInput)
			if command == "add" {
				todoList = append(todoList, messageOfTodo)
			} else if command == "done" {
				updateTodoList(messageOfTodo)
			}
			output := "Current Todos: \n"
			for _, todo := range todoList {
				output += "\n - " + todo + "\n"
			}

			output += "\n----------------------------------------"

			message = []byte(output)
			err = connection.WriteMessage(messageType, message)
			notifyAllConnections(connections)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	port := os.Getenv("PORT")

	http.ListenAndServe(":"+port, nil)
}
