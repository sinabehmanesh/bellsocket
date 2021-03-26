package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//wsHandler for websocket started pack :D
func WsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client CONNECTED!")
	err = ws.WriteMessage(1, []byte("Good day mate!"))
	if err != nil {
		log.Fatal(err)
	}

	WsReader(ws)
}

//WsReader to read client data
func WsReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Fatal(err)
			return
		}
	}
}

func main() {
	fmt.Println("hi sina")
	//fileserver
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	fmt.Println("file server UP!")

	//http.HandleFunc("/client", WsClient)

	http.HandleFunc("/bell", WsHandler)
	fmt.Println("WsHandler defiend!")
	//server on port 8080
	fmt.Println("Starting webService[:8080].....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
