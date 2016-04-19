package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type conn struct {
	c  *websocket.Conn
	ch chan bool
}

var connections = make(map[*conn]bool)
var connectionsLock = &sync.Mutex{}

func handleListen(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Error on eastablish connection: %s\n", err)
		return
	}

	myConn := &conn{
		c:  c,
		ch: make(chan bool),
	}

	connectionsLock.Lock()
	connections[myConn] = true
	connectionsLock.Unlock()

	go func() {
		defer func() {
			c.Close()
			connectionsLock.Lock()
			defer connectionsLock.Unlock()

			delete(connections, myConn)
			log.Println("cleanup")
		}()

		for range myConn.ch {
			var resp struct {
				N int `json:"n"`
			}

			resp.N = p.count()

			if err := c.WriteJSON(resp); err != nil {
				log.Printf("disconnect: %s", err)
				return
			}
		}
	}()
}
