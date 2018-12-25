package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{}

type ClientConnection struct {
	UserId      string
	ChannelName string
}

type Notification struct {
	Message string
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	var msg ClientConnection
	err = ws.ReadJSON(&msg)
	if err != nil {
		log.Printf("error: %v", err)
	}
	clients[msg.ChannelName] = make(map[string]*websocket.Conn)
	clients[msg.ChannelName][msg.UserId] = ws
}

func handleNewNotification(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	channelName := r.FormValue("channel_name")
	fmt.Println(userId)
	fmt.Println(channelName)
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/notify", handleNewNotification)
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
