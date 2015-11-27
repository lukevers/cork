/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package session

import (
	"github.com/gorilla/sessions"
	"github.com/srinathgs/mysqlstore"
)

type Session struct {
	Store string
	Secret string
	Stores map[string]sessions.Store
}

func (s *Session) New() {
	s.Stores = make(map[string]sessions.Store)

	switch s.Store {
		case "cookie":
			s.Stores["cookie"] = sessions.NewCookieStore([]byte(s.Secret))
		case "mysql":
			// TODO: THIS IS KINDA JUST A PLACEHOLDER, HAVEN'T TESTED IT OUT YET.
			mystore, _ := mysqlstore.NewMySQLStore("conn", "table", "/", 3600, []byte(s.Secret))
			defer mystore.Close()
			s.Stores["mysql"] = mystore
		case "sqlite":
			// TODO
	}
}

func (s *Session) GetStore() sessions.Store {
	return s.Stores[s.Store]
}
