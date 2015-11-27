/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"cork/session"
)

type HttpServer struct {
	Port int
	Host string
	Session *session.Session
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	store *sessions.Store
	Connections map[string]*websocket.Conn = make(map[string]*websocket.Conn)
)

type Message struct {
	Message string
}

func (s *HttpServer) Run() {
	s.Session.New()
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
	store := Cork.Http.Session.GetStore()
	session, _ := store.Get(r, "session-name")
	if session.IsNew {
		id, err := random(32)
		if err == nil {
			session.Values["Id"] = id
			session.Save(r, w)
		}
	}

	templates := template.Must(template.ParseGlob("resources/html/*.html"))
	templates.ExecuteTemplate(w, "index.html", struct {
		Id string
	}{
		Id: session.Values["Id"].(string),
	})
}

func ws(w http.ResponseWriter, r *http.Request) {
	store := Cork.Http.Session.GetStore()
	session, _ := store.Get(r, "session-name")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	Connections[session.Values["Id"].(string)] = conn

	for {
		var i map[string]string
		err := conn.ReadJSON(&i)
		if err != nil {
			break
		}

		log.Println(i)
		Keys[session.Values["Id"].(string)] = i["Value"]
	}
}

func random(size int) (string, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes), err
}
