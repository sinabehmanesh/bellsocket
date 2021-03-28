package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//Message object definition
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

//declare clients
var Clients = make(map[*websocket.Conn]bool)

//declare channel
var Broadcast = make(chan Message)

//declare incoming message srtucture
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//wsHandler for websocket started pack :D
func WsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new message!")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	//Upgrade GET connection!
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	//add client
	Clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Fatal(err)
			delete(Clients, ws)
			break
		}
		fmt.Println("data is ", msg)
		Broadcast <- msg
	}
	//	WsReader(ws)
}

//HandleMessage sends message data to all the clients!
func HandleMessage(wg *sync.WaitGroup) {
	for {
		msg := <-Broadcast
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Fatal(err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
	wg.Done()

}

//WsReader to read client data
func WsReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Client message: ")
		fmt.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Fatal(err)
			return
		}
	}
}

//SetupRoutes for handlers and URLS
func SetupRoutes(wg *sync.WaitGroup) {

	//fileserver
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	fmt.Println("file server UP!")

	//Handler html socket
	http.HandleFunc("/ws", WsHandler)

	//Done with the wait Group
	fmt.Println("Routes are SETTED!")
	wg.Done()
}

func main() {

	//add waitgroup (usage on SetupRoutes)
	var wg sync.WaitGroup

	wg.Add(1)
	go SetupRoutes(&wg)
	wg.Wait()

	//http.HandleFunc("/client", WsClient)
	fmt.Println("first goroutine is done!!")
	//wg.Add(1)
	go HandleMessage(&wg)
	wg.Wait()

	fmt.Println("WsHandler defiend!")
	//server on port 8080
	fmt.Println("Starting webService[:8000].....")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
