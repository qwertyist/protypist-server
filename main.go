package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/qwertyist/protypist/db"
	"github.com/qwertyist/protypist/handler"
	"github.com/qwertyist/protypist/session"
)

type config struct {
	port string
}

var cfg = config{port: ":4000"}

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	service := session.NewService(db)
	handler := handler.NewHandler(service)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", apiHelper).Methods("GET")
	http.Handle("/", accessControl(r))
	handler.Endpoints(r)
	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port:", cfg.port)
		errs <- http.ListenAndServe(cfg.port, nil)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("terminated: %s", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		///origin := req.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Id-Token")
		//w.Header().Set("Access-Control-Allow-Origin", "*")

		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func apiHelper(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the API helper speaking, available endpoints/methods are:"))
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This function is not implemented"))
}
