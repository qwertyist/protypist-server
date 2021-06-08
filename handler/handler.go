package handler

import (
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
	r.HandleFunc("/session", h.GetSession).Methods("GET")
	r.HandleFunc("/recv", h.RecvTX).Methods("GET")
}

func (h *handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Created session"))
}

func (h *handler) GetSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get session"))
}

func (h *handler) RecvTX(w http.ResponseWriter, r *http.Request) {

}
