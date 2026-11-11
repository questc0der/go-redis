package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, store *Store) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := readCommand(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("handleConnection: %v", err)
			}
			return
		}
		if len(command) == 0 {
			continue
		}

		h, ok := commands[strings.ToUpper(command[0])]
		if !ok {
			conn.Write(errorReply("unknown command '" + command[0] + "'"))
			continue
		}
		conn.Write(h(store, command))
	}
}