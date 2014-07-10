package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

var scripts map[string]string

const reloadScript = `

<script type='text/javascript'>
    var livereloadWebSocket = new WebSocket("ws://localhost:3000/reload/");
    livereloadWebSocket.onmessage = function(msg) {
        livereloadWebSocket.close();
        window.location.reload(true);
    };

    livereloadWebSocket.onopen = function(x) { console.log('[ws] Connection opened', new Date()); };
    livereloadWebSocket.onclose = function() { console.log('[ws] closing'); };
    livereloadWebSocket.onerror = function(err) { console.log('[ws] error', err); };

    window.define = function(){};
</script>
`

func handleIndex(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	file = file
	fmt.Fprintf(w, reloadScript+string(file))
}

func handleAssets(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/assets/"):]

	if file == "app.js" {
		w.Header().Set("Content Type", "text/javascript")

		for _, script := range scripts {
			fmt.Fprint(w, script+"\n\n")
		}
	}
}

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

func createClient(ws *websocket.Conn) {
	log.Println(Color("[ws]", "yellow"), "Client connected")

	client := Client{clientId, ws, make(chan bool)}
	clients[client.id] = &client
	clientId++

	client.Listen()
}

func StartServer(port string, c chan map[string]string) {
	log.Println(Color("[server]", "green"), "Starting server on port", port)

	clients = make(map[int]*Client)

	http.Handle("/reload/", websocket.Handler(createClient))

	http.HandleFunc("/assets/", handleAssets)
	http.HandleFunc("/", handleIndex)

	go func() {
		for {
			select {
			case js := <-c:
				log.Println(Color("[server]", "green"), "updated scripts")
				scripts = js

				for _, client := range clients {
					client.reloadCh <- true
				}

			}
		}
	}()

	http.ListenAndServe(":"+port, nil)

}
