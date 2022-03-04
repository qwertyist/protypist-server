package handler

import (
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
	w.Write(h.service.ReadBuf(uuid))
}
