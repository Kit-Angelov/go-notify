package main

import (
	// "bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var clients = make(map[string][]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query()["channel"][0]
	println(channel)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[channel] = append(clients[channel], ws)
	fmt.Println(clients)

	for {
		mt, msg, _ := ws.ReadMessage()
		println(string(msg))
		println("mt", mt)
		for _, client := range clients[channel] {
			client.WriteMessage(mt, msg)
		}
	}
}

func handleNewNotification(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	channelName := r.FormValue("channel_name")
	fmt.Println(userId)
	fmt.Println(channelName)
}

func main() {
	fs := http.FileServer(http.Dir(""))
	http.Handle("/", fs)
	http.HandleFunc("/ws/", handleConnections)
	http.HandleFunc("/broadcast/", handleNewNotification)
	log.Println("http server started on :8000")
	err := http.ListenAndServe("192.168.0.105:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
