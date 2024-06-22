package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func handleConnection(c net.Conn) {
	request, err := http.ReadRequest(bufio.NewReader(c))
	if err != nil {
		fmt.Println("Failed to read request", err.Error())
	}

	switch {
	case request.URL.Path == "/":
		_, err = c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	case strings.HasPrefix(request.URL.Path, "/files/"):
		var dir string
		if os.Args[1] == "--directory" {
			dir = os.Args[2]
		} else {
			fmt.Fprintf(os.Stderr, "`--directory` flag is missing")
			os.Exit(1)
		}

		filename := request.URL.Path[7:]
		content, err := os.ReadFile(fmt.Sprintf("%s%s", dir, filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			_, err = c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			break
		}
		data := string(content)
		_, err = c.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)))

	case request.URL.Path == "/user-agent":
		_, err = c.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request.UserAgent()), request.UserAgent())))

	case strings.HasPrefix(request.URL.Path, "/echo/"):
		str := request.URL.Path[6:]
		contentLength := len(str)

		_, err = c.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, str)))

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
