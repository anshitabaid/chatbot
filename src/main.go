package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make (map[*websocket.Conn]bool) //connected clients
var broadcast = make (chan Message)

//configure the upgrader
var upgrader = websocket.Upgrader{}

//message object
type Message struct {
	Email		string	`json:"email"`
	Username	string	`json: username`
	Message		string	`json:"message"`
}

func main (){
	//simple file server
	fs := http.FileServer (http.Dir("../public"))
	http.Handle ("/", fs)
	//configure web sockets
	http.HandleFunc ("/ws", handleConnections)
	go handleMessages ()
	log.Println ("http server started on :8000")
	err := http.ListenAndServe (":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections (w http.ResponseWriter, r *http.Request){
	//upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade (w, r, nil)
	if err != nil {
		log.Fatal (err)
	}
	//close connection when function returns
	defer ws.Close ()

	//register new client
	clients [ws] = true

	for {
		var msg Message
		//read message as JSON and map it to Message object
		err := ws.ReadJSON (&msg)
		if err != nil {
			log.Printf ("error: %v", err)
			delete (clients, ws)
			break
		}
		//send new message to broadcast
		broadcast <- msg
	}
}
func handleMessages (){
	for {
		//grab next message from broadcast channel
		msg := <-broadcast
		//send to every client currently connected
		for client := range clients {
			err := client.WriteJSON (msg)
			if err!= nil {
				log.Printf ("error: %v", err)
				client.Close ()
				delete (clients, client)
			}
		}
	}
}


