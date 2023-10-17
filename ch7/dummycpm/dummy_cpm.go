package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listen, err := net.Listen("unix", "/tmp/cpm.sock")

	if err != nil {
		log.Fatalln("Socket listen failed:", err)
	}

	handleSigTerm(listen)

	pbListener := NewProtobufListener()
	for {
		conn, err := listen.Accept()
		if err == nil {
			go pbListener.Handle(conn)
			log.Println("Started CPM mock")
		} else {
			log.Println("Error waiting for incoming connections:", err)
			var err = listen.Close()
			if err != nil {
				log.Println("Error closing", err)
			}
		}
	}
}

func handleSigTerm(listen net.Listener) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopped CPM mock")
		var err = os.Remove("/tmp/cpm.sock")
		if err != nil {
			log.Println("Error removing", err)
		}
		log.Println("Cleaned up")
		os.Exit(0)
	}()
}
