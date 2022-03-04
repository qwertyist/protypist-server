package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qwertyist/protypist/session"
)

type Handler interface {
	CreateSession(w http.ResponseWriter, r *http.Request)
	GetSession(w http.ResponseWriter, r *http.Request)
	RecvTX(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service *session.Service
}

func NewHandler(service *session.Service) *handler {
	h := &handler{
		service: service,
	}
	return h
}

func (h *handler) Endpoints(r *mux.Router) {
	r.HandleFunc("/session", h.CreateSession).Methods("POST")
	r.HandleFunc("/session/{uuid}", h.GetSession).Methods("GET")
	r.HandleFunc("/ws", h.wsInterpreterHandler)
	r.HandleFunc("/listen/{uuid}", h.wsListenerHandler)

}

func (h *handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	uuid := h.service.CreateSession("")
	log.Printf("Created sesssion: %s", uuid)
	w.Write([]byte(uuid))
}

func (h *handler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	s := h.service.GetSession(uuid)

	sess := struct {
		ID        string
		Connected int
		Clients   []*session.Client
		Buf       []byte
	}{
		ID:        uuid,
		Connected: len(s.Clients),
		Clients:   s.Clients,
		Buf:       s.Buf,
	}

	bytes, err := json.MarshalIndent(sess, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Marshall failed"))
	}

	w.Write(bytes)
}
