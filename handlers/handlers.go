package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

type wsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// UpgradeConnection is used to upgrade the connection to a websocket connection
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint ")

	var response wsJsonResponse
	response.Message = `<em><small>Connected to server <small></em>`

	// When a client connects to the websocket endpoint, we create a new WebSocketConnection
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	// Start listening for messages from the client
	go ListenForWs(&conn)
}

func broadCastToAll(response wsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println(err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error: ", fmt.Sprintf("%s", r))
		}
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// Show the message to all clients

func ListenToWsChannel() {
	var response wsJsonResponse
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

// renderPage is used to render a page using the jet templating engine
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
