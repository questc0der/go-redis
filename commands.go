package main

import (
	"strconv"
	"strings"
	"time"
)

type handler func(store *Store, args []string) []byte

var commands = map[string]handler{
	"PING":    cmdPing,
	"ECHO":    cmdEcho,
	"SET":     cmdSet,
	"GET":     cmdGet,
	"COMMAND": cmdCommand,
}

func cmdPing(store *Store, args []string) []byte {
	return simpleString("PONG")
}

func cmdCommand(store *Store, args []string) []byte {
	return emptyArray()
}

func cmdEcho(store *Store, args []string) []byte {
	if len(args) != 2 {
		return errorReply("wrong number of arguments for 'echo' command")
	}
	return bulkString(args[1])
}

func cmdSet(store *Store, args []string) []byte {
	if len(args) != 3 && len(args) != 5 {
		return errorReply("wrong number of arguments for 'set' command")
	}

	var expiresAt time.Time
	if len(args) == 5 {
		if strings.ToUpper(args[3]) != "EX" {
			return errorReply("syntax error")
		}
		seconds, err := strconv.Atoi(args[4])
		if err != nil || seconds <= 0 {
			return errorReply("value is not an integer or out of range")
		}
		expiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
	}

	store.Set(args[1], args[2], expiresAt)
	return simpleString("OK")
}

func cmdGet(store *Store, args []string) []byte {
	if len(args) != 2 {
		return errorReply("wrong number of arguments for 'get' command")
	}
	value, ok := store.Get(args[1])
	if !ok {
		return nullBulk()
	}
	return bulkString(value)
}