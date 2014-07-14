package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"code.google.com/p/go.net/websocket"
)

var (
	proxy     *httputil.ReverseProxy
	port      string
	scripts   map[string]string
	vendorJs  []string
	vendorCss []string
)

func reloadScript() string {
	return `
<script type='text/javascript'>
    var livereloadWebSocket = new WebSocket("ws://localhost:` + port + `/reload/");
    livereloadWebSocket.onmessage = function(msg) {
        livereloadWebSocket.close();
        window.location.reload(true);
    };

    livereloadWebSocket.onopen = function(x) { console.log('[ws] Connection opened', new Date()); };
    livereloadWebSocket.onclose = function() { console.log('[ws] closing'); };
    livereloadWebSocket.onerror = function(err) { console.log('[ws] error', err); };
</script>
`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(proxy)
	if r.URL.Path[0:] != "/" && proxy != nil {
		proxy.ServeHTTP(w, r)
		return
	}

	file, err := ioutil.ReadFile("app/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, reloadScript()+string(file))
}

func handleAssets(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get current working directory")
	}

	file := r.URL.Path[len("/assets/"):]

	switch file {
	case path.Base(wd) + ".js": // TODO - allow override within config file
		w.Header().Set("Content Type", "text/javascript")

		for _, script := range scripts {
			fmt.Fprint(w, script+"\n")
		}

	case "vendor.js":
		w.Header().Set("Content Type", "text/javascript")

		for _, vendor := range vendorJs {
			fmt.Fprint(w, vendor+"\n\n")
		}

	case "app.css":
		w.Header().Set("Content Type", "text/css")
		http.ServeFile(w, r, "app/styles/app.css")

	case "vendor.css":
		w.Header().Set("Content Type", "text/css")
		for _, css := range vendorCss {
			fmt.Fprint(w, css+"\n")
		}
	}
}

func StartServer(portt string, prox string, fileC chan *File) {
	log.Println(Color("[server]", "green"), "Starting server on port", port)

	port = portt

	if prox != "" {
		p, err := url.Parse(prox)
		if err != nil {
			log.Fatal("Error parsing proxy URL", err)
		}

		proxy = httputil.NewSingleHostReverseProxy(p)
	}

	scripts = make(map[string]string)
	vendorJs, _ = getVendors("js")
	vendorCss, _ = getVendors("css")

	go listenForFiles(fileC)

	http.Handle("/reload/", websocket.Handler(CreateClient))
	http.HandleFunc("/assets/", handleAssets)
	http.HandleFunc("/", handleIndex)

	http.ListenAndServe(":"+port, nil)
}

func getVendors(kind string) (scripts []string, err error) {
	dir := "vendor/"

	for vendor, path := range Config.Vendors[kind] {
		file, err := ioutil.ReadFile(filepath.Join(dir, path))
		if err != nil {
			log.Fatal("error reading vendor file:", vendor, err)
		}

		scripts = append(scripts, string(file))
	}
	return
}

func listenForFiles(fileC chan *File) {
	for {
		select {
		case f := <-fileC:
			// TODO - handle by file type
			if f.IsEmpty() {
				delete(scripts, f.Path)
			} else {
				scripts[f.Path] = string(f.Content)
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
