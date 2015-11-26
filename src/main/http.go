/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/codegangsta/negroni"
	"html/template"
)

type HttpServer struct {
	Port int
	Host string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var Connections map[string]*websocket.Conn = make(map[string]*websocket.Conn)

type Message struct {
	Message string
}

func (s *HttpServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", root)
	router.HandleFunc("/ws", ws)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}

func UserHasBeenAuthorized(id string) {
	log.Println("Telling the user they have been authorized...")
	for i, conn := range Connections {
		if id == i {
			conn.WriteJSON(&Message{
				Message: "SUCCESS",
			})
		}
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseGlob("resources/html/*.html"))
	templates.ExecuteTemplate(w, "index.html", struct{
		Id string
	}{
		Id: "test",
	})
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	Connections["test"] = conn

	for {
		var i map[string]string
		err := conn.ReadJSON(&i)
		if err != nil {
			break
		}

		log.Println(i)
		Keys[i["Id"]] = i["Value"]
	}
}
