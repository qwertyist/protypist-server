package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *handler) wsInterpreterHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.Header.Get("session")
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	conns := h.service.GetSession(uuid)
	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		h.service.WriteBuf(uuid, msg)
		for _, c := range conns {
			c.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (h *handler) wsListenerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
	h.service.JoinSession(uuid, conn)
}
