package main

import (
	"flag"
	"html/template"
	_ "image/png"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	_ "github.com/boombuler/barcode"
	_ "github.com/boombuler/barcode/qr"

	"github.com/jmcvetta/randutil"

	_ "github.com/scritch007/Racer/types"
	"github.com/scritch007/go-tools"
)

var addr = flag.String("addr", ":8080", "http service address")
var debug = flag.Bool("debug", false, "Turn into debug mode")

var serverConnections map[string]*websocket.Conn

func init() {
	serverConnections = make(map[string]*websocket.Conn)
}

var upgrader = websocket.Upgrader{} // use default options
func serverWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tools.LOG_ERROR.Print("upgrade:", err)
		return
	}

	// Notify the others that we have a new connection pending
	for _, conn := range serverConnections {
		tools.LOG_DEBUG.Println("New message to send")
		conn.WriteMessage(1, []byte(`{"id":"`+id+`"}`))
	}

	tools.LOG_DEBUG.Println("Entering")

	serverConnections[id] = c

	defer func() {
		c.Close()
		delete(serverConnections, id)
	}()

	for {
		mt, m, err := c.ReadMessage()
		if err != nil {
			tools.LOG_ERROR.Println("read:", err)
			break
		}
		//Forward to everybody except us
		for cId, conn := range serverConnections {
			if cId == id {
				continue
			}
			conn.WriteMessage(mt, m)
		}

	}
}
func serveFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filePath := vars["file"]
	http.ServeFile(w, r, "./html/"+filePath)
}
func home(w http.ResponseWriter, r *http.Request) {
	var id string

	//Create unique RandomUser
	id, _ = randutil.AlphaString(20)

	tmpl, err := template.ParseFiles("./html/index.html")
	if err != nil {
		panic(err)
	}
	var socketProto = "ws"
	if strings.Contains(r.Proto, "HTTPS") {
		socketProto = "wss"
	}
	var first bool
	if 0 == len(serverConnections) {
		first = true
	} else {
		first = false
	}
	values := struct {
		WSocketURL string
		Id         string
		First      bool
	}{
		WSocketURL: socketProto + "://" + r.Host + "/ws/" + id,
		Id:         id,
		First:      first,
	}
	tmpl.Execute(w, values)
}

func main() {
	tools.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	flag.Parse()
	tools.LOG_ERROR.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/static/{file:.*}", serveFile)
	r.HandleFunc("/ws/{id}", serverWS)
	r.HandleFunc("/", home)
	http.Handle("/", r)
	tools.LOG_ERROR.Fatal(http.ListenAndServe(*addr, nil))
}
