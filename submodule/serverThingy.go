package submodule

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type server struct {
	Name           string
	MessageChannel chan string
	ControlChannel chan string
	Port           string
	Host           string
	Router         *mux.Router
}

func New(name, host string, port int, messageChannel, controlChannel chan string) *server {
	s := server{
		Name:           name,
		Host:           host,
		Port:           ":" + string(port),
		MessageChannel: messageChannel,
		ControlChannel: controlChannel,
	}

	r := mux.NewRouter()
	r.HandleFunc("/close/", s.CloseHandler)
	r.HandleFunc("/msg/{foo}", s.DefaultHandler)
	r.NotFoundHandler = http.HandlerFunc(s.NotFoundHandler)
	s.Router = r

	return &s
}

func (s *server) Connect() {
	srv := &http.Server{
		Handler:      s.Router,
		Addr:         s.Host + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go s.handleIPC()
	log.Fatal(srv.ListenAndServe())
}

func (s *server) handleIPC() {
	for {
		select {
		case msg := <-s.MessageChannel:
			fmt.Println("Received (" + s.Name + "): " + msg)

		case msg := <-s.ControlChannel:
			fmt.Println(msg)
		}
	}
}

func (s *server) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	s.MessageChannel <- s.Name + " " + mux.Vars(r)["foo"]
}

func (s *server) CloseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Close")
	s.ControlChannel <- "close"
}

func (s *server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Nicht gefunden, dude")
}
