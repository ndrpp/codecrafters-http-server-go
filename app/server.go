package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func handleConnection(c net.Conn) {
	request, err := http.ReadRequest(bufio.NewReader(c))
	if err != nil {
		fmt.Println("Failed to read request", err.Error())
	}

	switch request.URL.Path {
	case "/":
		_, err = c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		_, err = c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	if err != nil {
		fmt.Println("Failed to write response to socket", err.Error())
	}

	c.Close()
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("TCP Server listening on port: ", 4221)

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			handleConnection(c)
		}(conn)
	}
}
