package main

import  (
	"net"
	"log"
	"fmt"
	"io"
	"bufio"
	"strconv"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var command []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("handleConnection: %v", err)
			}
			return
		}
		if line[0] == '*'{
		strTrm := strings.TrimSpace(line[1:])
		num,_ := strconv.Atoi(strTrm)
		for i := 0; i < num; i++ {
		line2, _ := reader.ReadString('\n')
			if line2[0] == '$'{
					line3, _ := reader.ReadString('\n')
					command = append(command, string(strings.TrimSpace(line3)))
			}
		}
	}	
		fmt.Printf("%q\n", line)
		if strings.ToUpper(command[0]) == "PING"{
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}

func main () {
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Printf("ERR: %v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Err: %v", err)
		}
		go handleConnection(conn)
	}
}
