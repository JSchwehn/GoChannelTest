package submodule

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type otherServer struct {
	Name           string
	MessageChannel chan string
	ControlChannel chan string
	Port           int
	Host           string
	Router         *mux.Router
}

func NewOtherServer(name, host string, port int, messageChannel, controlChannel chan string) *otherServer {
	s := otherServer{
		Name:           name,
		Host:           host,
		Port:           port,
		MessageChannel: messageChannel,
		ControlChannel: controlChannel,
	}

	r := mux.NewRouter()
	r.HandleFunc("/close/", s.CloseHandler)
	r.HandleFunc("/msg/{foo}", s.DefaultHandler)

	s.Router = r

	return &s
}

func (s *otherServer) Connect() {
	srv := &http.Server{
		Handler:      s.Router,
		Addr:         s.Host + ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go s.handleIPC()
	log.Fatal(srv.ListenAndServe())
}

func (s *otherServer) handleIPC() {
	for {
		select {
		case msg := <-s.MessageChannel:
			fmt.Println("Received (" + s.Name + "): " + msg)

		case msg := <-s.ControlChannel:
			fmt.Println(msg)
		}
	}
}

func (s *otherServer) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	s.MessageChannel <- s.Name + " " + mux.Vars(r)["foo"]
}

func (s *otherServer) CloseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Close")
	s.ControlChannel <- "close"
}
