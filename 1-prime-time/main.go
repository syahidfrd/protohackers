package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"
	"net"
)

type request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("line: %s\n", line)

		var req request
		err := json.Unmarshal([]byte(line), &req)
		if err != nil {
			log.Printf("error malformed request: %v", err)
			conn.Write([]byte("💀"))
			return
		}

		if req.Method != "isPrime" || req.Number == nil {
			log.Printf("error malformed request: invalid method or number")
			conn.Write([]byte("💀"))
			return
		}

		resp := response{
			Method: "isPrime",
			Prime:  isPrime(*req.Number),
		}

		js, _ := json.Marshal(resp)
		conn.Write([]byte(append(js, '\n')))
	}

	log.Println("connection closed")
}

func isPrime(n float64) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}
