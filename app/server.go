package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

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
			fmt.Println("Connection from: ", c.RemoteAddr().String())

			buffer := make([]byte, 1024)
			_, err := c.Read(buffer)
			if err != nil {
				fmt.Println("Failed to read request", err.Error())
			}
			data := strings.Split(string(buffer), "\r\n")
			path := data[0][4 : len(data[0])-9]

			switch path {
			case string("/"):
				_, err = c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
				if err != nil {
					fmt.Println("Failed to write response to socket", err.Error())
				}
			default:
				_, err = c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				if err != nil {
					fmt.Println("Failed to write response to socket", err.Error())
				}
			}
			c.Close()
		}(conn)
	}
}
