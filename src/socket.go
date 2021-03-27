package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
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
	//Upgrade GET connection!
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	log.Println("\nClient ", r.RemoteAddr, " CONNECTED!")
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
	http.HandleFunc("/bell", WsHandler)

	//Done with the wait Group
	wg.Done()
}

type SocketData struct {
	Text string `name:"text"`
}

func main() {
	RedisPort := os.Getenv("PORT")
	fmt.Println(RedisPort)
	//add waitgroup (usage on SetupRoutes)
	var wg sync.WaitGroup

	wg.Add(1)
	go SetupRoutes(&wg)
	wg.Wait()

	//http.HandleFunc("/client", WsClient)

	fmt.Println("WsHandler defiend!")
	//server on port 8080
	fmt.Println("Starting webService[:8080].....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
