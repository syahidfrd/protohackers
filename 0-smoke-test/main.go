package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

	defer listener.Close()
	log.Println("server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}

		log.Println("client connected")
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			log.Println("client finished sending data")
			break
		}

		if err != nil {
			log.Printf("error reading from client: %v", err)
			return
		}

		log.Printf("received %d bytes: %x", n, buffer[:n])
		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Printf("error writing to client: %v", err)
			return
		}
	}
	log.Println("connection closed")
}
