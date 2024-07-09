package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type ResponseStatus string

const (
	OK        = "200 OK"
	NOT_FOUND = "404 Not Found"
)

type ServerRequest struct {
	Method  string
	Uri     string
	Version string
}

func handleRequest(data []byte) ServerRequest {
	parts := strings.Split(string(data), " ")
	return ServerRequest{
		Method:  parts[0],
		Uri:     parts[1],
		Version: parts[2],
	}
}

func writeResponse(conn net.Conn, status ResponseStatus) (int, error) {
	return conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %s\r\n\r\n", status)))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	data := make([]byte, 1024)
	_, err = connection.Read(data)
	if err != nil {
		log.Fatalln("Error reading from connection", err)
	}

	if request := handleRequest(data); request.Uri != "/" {
		_, err = writeResponse(connection, NOT_FOUND)
	} else {
		_, err = writeResponse(connection, OK)
	}

	if err != nil {
		log.Fatalln("Couldn't respond the connection!")
	}
}
