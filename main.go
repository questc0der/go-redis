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
	"time"
)

type Expire struct {
	Value     string
	ExpiresAt time.Time // zero value = no expiry
}

type Args struct {
	args map[string]Expire
	rmw  sync.RWMutex
}

func sweeper(store *Args) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		store.rmw.Lock()
		for key, entry := range store.args {
			if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
				delete(store.args, key)
			}
		}
		store.rmw.Unlock()
	}
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
			if len(command) != 3 && len(command) != 5 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}

			entry := Expire{Value: command[2]}

			if len(command) == 5 {
				if strings.ToUpper(command[3]) != "EX" {
					conn.Write([]byte("-ERR syntax error\r\n"))
					continue
				}
				seconds, err := strconv.Atoi(command[4])
				if err != nil || seconds <= 0 {
					conn.Write([]byte("-ERR value is not an integer or out of range\r\n"))
					continue
				}
				entry.ExpiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
			}

			args.rmw.Lock()
			args.args[command[1]] = entry
			args.rmw.Unlock()
			conn.Write([]byte("+OK\r\n"))
		case "GET":
			if len(command) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}

			args.rmw.RLock()
			entry, ok := args.args[command[1]]
			args.rmw.RUnlock()

			if !ok || (!entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt)) {
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			conn.Write(fmt.Appendf(nil, "$%d\r\n%s\r\n", len(entry.Value), entry.Value))
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
	args := &Args{args: make(map[string]Expire)}
	go sweeper(args)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Err: %v", err)
		}
		go handleConnection(conn, args)
	}
}
