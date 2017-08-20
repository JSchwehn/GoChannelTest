package main

import (
	"github.com/jschwehn/channelTest/submodule"
)

func main() {
	msgCh := make(chan string)
	ctlCh := make(chan string)

	defer close(msgCh)
	defer close(ctlCh)

	srv := submodule.New("ServerThingy", "localhost", 8080, msgCh, ctlCh)
	otherSrv := submodule.NewOtherServer("OtherServer", "localhost", 8081, msgCh, ctlCh)

	go srv.Connect()
	go otherSrv.Connect()

	for {
		select {
		case command := <-ctlCh:
			if command == "close" {
				return
			}
		}
	}

}
