package lib

import (
	"log"

	"code.google.com/p/go.net/websocket"
)

// TODO - remove clients when browser is manually refreshed
type Client struct {
	id       int
	conn     *websocket.Conn
	reloadCh chan bool
}

func (c *Client) Listen() {
	for {
		select {
		case <-c.reloadCh:
			c.conn.Write([]byte("RELOAD"))
			c.conn.Close()
			delete(clients, c.id)
		}
	}
}

var clientId = 0
var clients map[int]*Client

func CreateClient(ws *websocket.Conn) {
	log.Println(Color("[ws]", "yellow"), "Client connected")

	client := Client{clientId, ws, make(chan bool)}
	clients[client.id] = &client
	clientId++

	client.Listen()
}

func init() {
	clients = make(map[int]*Client)
}
