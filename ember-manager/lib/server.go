package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"code.google.com/p/go.net/websocket"
)

var scripts map[string]string
var vendorScripts []string

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
</script>
`

func handleIndex(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, reloadScript+string(file))
}

func handleAssets(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/assets/"):]

	switch file {
	case "app.js":
		w.Header().Set("Content Type", "text/javascript")

		for _, script := range scripts {
			fmt.Fprint(w, script+"\n")
		}

	case "vendor.js":
		w.Header().Set("Content Type", "text/javascript")

		for _, vendor := range vendorScripts {
			fmt.Fprint(w, vendor+"\n\n")
		}

	case "app.css":
		w.Header().Set("Content Type", "text/css")
		http.ServeFile(w, r, "app/styles/app.css")
	}
}

func StartServer(port string, fileC chan *File) {
	log.Println(Color("[server]", "green"), "Starting server on port", port)

	scripts = make(map[string]string)
	vendorScripts, _ = getVendorJS()
	go listenForFiles(fileC)

	http.Handle("/reload/", websocket.Handler(CreateClient))
	http.HandleFunc("/assets/", handleAssets)
	http.HandleFunc("/", handleIndex)

	http.ListenAndServe(":"+port, nil)

}

func getVendorJS() (vendorScripts []string, err error) {
	dir := "vendor/"

	for vendor, path := range Config.Vendors {
		file, err := ioutil.ReadFile(filepath.Join(dir, path))
		if err != nil {
			log.Fatal("error reading vendor file:", vendor, err)
		}

		vendorScripts = append(vendorScripts, string(file))
	}
	return
}

func listenForFiles(fileC chan *File) {
	for {
		select {
		case f := <-fileC:
			if len(f.Content) > 0 {

				scripts[f.Path] = string(f.Content)
			} else {
				delete(scripts, f.Path)
			}
			reloadAllClients()
		}
	}
}

func reloadAllClients() {
	log.Println(Color("[server]", "green"), "reloading clients")

	for _, client := range clients {
		client.reloadCh <- true
	}
}
