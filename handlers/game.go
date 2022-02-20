package handlers

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ResponseType string

const (
	Join        ResponseType = "join"
	Move        ResponseType = "move"
	NewPlayer   ResponseType = "newPlayer"
	ColorChange ResponseType = "colorChange"
	Leave       ResponseType = "leave"
)

type OnlinePlayer struct {
	Id         string `json:"id"`
	Position   []int  `json:"position"`
	Color      string `json:"color"`
	Connection *websocket.Conn
}

type Message struct {
	Type     string `json:"type"`
	IdPlayer string `json:"idPlayer"`
}

type JoinMessage struct {
	Type        string `json:"type"`
	IdPlayer    string `json:"idPlayer"`
	PlayerColor string `json:"color"`
	Position    []int  `json:"position"`
}

type MoveMessage struct {
	Type     string `json:"type"`
	IdPlayer string `json:"idPlayer"`
	Position []int  `json:"position"`
}

type ColorChangeMessage struct {
	Type     string `json:"type"`
	IdPlayer string `json:"idPlayer"`
	Color    string `json:"color"`
}

var onlinePlayers []OnlinePlayer

func HandleMessages(connection *websocket.Conn) error {
	var baseMessage Message
	_, message, err := connection.ReadMessage()
	if err != nil {

		log.Println(connection.RemoteAddr().String())
		log.Println(err)
		return err
	}

	err = json.Unmarshal(message, &baseMessage)

	if err != nil {
		log.Println(err)
		return err
	}
	switch ResponseType(baseMessage.Type) {
	case Join:
		var joinMessage JoinMessage
		err = json.Unmarshal(message, &joinMessage)
		if err != nil {
			log.Println(err)
			return err
		}

		newPlayer := OnlinePlayer{
			Id:         joinMessage.IdPlayer,
			Position:   joinMessage.Position,
			Color:      joinMessage.PlayerColor,
			Connection: connection,
		}

		createPlayer(newPlayer)
		notifyAllConnections(Join, onlinePlayers)
		notifyAllConnections(NewPlayer, newPlayer)
		log.Println(onlinePlayers)

		log.Println(newPlayer)
	case Move:
		var moveMessage MoveMessage
		err = json.Unmarshal(message, &moveMessage)
		if err != nil {
			log.Println(err)
			return err
		}

		for index := range onlinePlayers {
			player := &onlinePlayers[index]
			if player.Id == moveMessage.IdPlayer {
				player.Position = moveMessage.Position
				notifyAllConnections(Move, player)
				break
			}
		}

	case ColorChange:
		var colorChangeMessage ColorChangeMessage
		err = json.Unmarshal(message, &colorChangeMessage)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		for index := range onlinePlayers {
			player := &onlinePlayers[index]
			if player.Id == colorChangeMessage.IdPlayer {
				player.Color = colorChangeMessage.Color
				notifyAllConnections(ColorChange, player)
				break
			}
		}
	}

	return nil
}

func HandleDisconnect(connection *websocket.Conn) {
	for index := range onlinePlayers {
		player := &onlinePlayers[index]
		if player.Connection == connection {
			// remove player
			onlinePlayers = append(onlinePlayers[:index], onlinePlayers[index+1:]...)

			notifyAllConnections(Leave, player)
			break
		}
	}
}

func createPlayer(newPlayer OnlinePlayer) {
	onlinePlayers = append(onlinePlayers, newPlayer)
}

type ResponseMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func notifyAllConnections(messageType ResponseType, message interface{}) {
	responseMessage := ResponseMessage{
		Type: string(messageType),
		Data: message,
	}
	for _, player := range onlinePlayers {
		player.Connection.WriteJSON(responseMessage)
	}
}
