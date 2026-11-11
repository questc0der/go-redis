package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const maxBulkLen = 512 * 1024 * 1024

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func readNumber(r *bufio.Reader, prefix byte) (int, error) {
	line, err := readLine(r)
	if err != nil {
		return 0, err
	}
	if len(line) == 0 || line[0] != prefix {
		return 0, fmt.Errorf("protocol error: expected %q, got %q", prefix, line)
	}
	return strconv.Atoi(line[1:])
}

func readCommand(r *bufio.Reader) ([]string, error) {
	n, err := readNumber(r, '*')
	if err != nil {
		return nil, err
	}
	if n < 0 || n > 1024 {
		return nil, fmt.Errorf("protocol error: invalid array length %d", n)
	}

	command := make([]string, 0, n)
	for i := 0; i < n; i++ {
		size, err := readNumber(r, '$')
		if err != nil {
			return nil, err
		}
		if size < 0 || size > maxBulkLen {
			return nil, fmt.Errorf("protocol error: invalid bulk length %d", size)
		}
		buf := make([]byte, size+2) // data + \r\n
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		if buf[size] != '\r' || buf[size+1] != '\n' {
			return nil, fmt.Errorf("protocol error: bulk string not terminated by CRLF")
		}
		command = append(command, string(buf[:size]))
	}
	return command, nil
}

func simpleString(s string) []byte { return []byte("+" + s + "\r\n") }
func errorReply(msg string) []byte { return []byte("-ERR " + msg + "\r\n") }
func nullBulk() []byte             { return []byte("$-1\r\n") }
func emptyArray() []byte           { return []byte("*0\r\n") }
func bulkString(s string) []byte {
	return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(s), s)
}