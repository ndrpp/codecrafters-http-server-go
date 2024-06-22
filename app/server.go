package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

			_, err := c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			if err != nil {
				fmt.Println("Failed to write response to socket", err.Error())
			}
			c.Close()
		}(conn)
	}
}
