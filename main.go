package main

import  (
	"net"
	"log"
	"fmt"
	"io"
	"bufio"
	"strconv"
	"strings"
	"sync"
)

type Args struct{
	args map[string]string
	rmw sync.RWMutex
}


func handleConnection(conn net.Conn, args *Args) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		var command []string

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
		if len(command[0]) == 0 {
			continue
		}
		switch strings.ToUpper(command[0]){
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			var response []byte
			if len(command) > 1 {
				response = fmt.Appendf(response, "$%d\r\n%s\r\n", len(command[1]), command[1])
				conn.Write(response)
			} else{
				conn.Write([]byte("-ERR Wrong number of argument for 'echo' command\r\n"))
			}
		case "COMMAND":
			conn.Write([]byte("+Empty\r\n"))
		case "SET":
			if len(command) != 3{
				conn.Write([]byte("-ERR wrong number of argument for 'set' command\r\n"))
				continue
			}else {
			args.rmw.Lock()
			args.args[command[1]] = command[2]
			args.rmw.Unlock()
			
			conn.Write([]byte("+OK\r\n"))
			}
		case "GET":
			if len(command) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			
			args.rmw.Lock()
			value, ok := args.args[command[1]]
			args.rmw.Unlock()

			if !ok{
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			response := fmt.Appendf(nil, "$%d\r\n%s\r\n", len(value), value)
			conn.Write(response)
		default:
			conn.Write([]byte("-Error unknown command\r\n"))
		}
		
	}
}

func main () {
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Printf("ERR: %v", err)
	}
	args := &Args{args: make(map[string]string)}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Err: %v", err)
		}
		go handleConnection(conn, args)
	}
}
