package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	store := NewStore()
	go store.sweep()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleConnection(conn, store)
	}
}