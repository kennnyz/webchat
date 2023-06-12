package handlers

import (
	"database/sql"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"github.com/kennnyz/webchat/models"
	"log"
	"net/http"
	"sort"
)

var wsChan = make(chan models.ClientNotifier)
var clients = make(map[models.WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// UpgradeConnection is used to upgrade the connection to a websocket connection
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Мы подключились к веб сокету и теперь слушаем его в горутине ListenForWs()
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint ")

	var response models.WsJsonResponse
	response.Message = `<em><small>Connected to server <small></em>`

	// When a client connects to the websocket endpoint, we create a new WebSocketConnection
	conn := models.WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	// Start listening for messages from the client
	go ListenForWs(&conn)
}

func broadCastToAll(response models.WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println(err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenForWs(conn *models.WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error: ", fmt.Sprintf("%s", r))
		}
	}()
	var notifier models.ClientNotifier
	for {
		err := conn.ReadJSON(&notifier)
		if err != nil {
			// ничего не делаем
		} else {
			notifier.Conn = *conn
			wsChan <- notifier
		}
	}
}

// Show the message to all clients
// should listen

// Слушаем события от пользователей и транслируем его другим

func ListenToWsChannel(db *sql.DB) {
	var response models.WsJsonResponse
	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			// get list of all users and send it back via broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadCastToAll(response)
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadCastToAll(response)
		case "broadcast":
			err := e.WriteActionToDB(db)
			if err != nil {
				log.Fatal(err)
			}
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadCastToAll(response)
		}
	}
}

func getUserList() []string {
	var users []string
	for _, v := range clients {
		if v != "" {
			users = append(users, v)
		}
	}
	sort.Strings(users)
	return users
}
