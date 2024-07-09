package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type ResponseStatus string

const (
	OK        ResponseStatus = "200 OK"
	NOT_FOUND ResponseStatus = "404 Not Found"
)

type ServerRequest struct {
	Method  string
	Uri     string
	Version string
}

type ServerResponse struct {
	status  ResponseStatus
	headers map[string]string
	body    string
}

func (sr *ServerResponse) SetHeader(header string, value string) *ServerResponse {
	sr.headers[header] = value
	return sr
}

func (sr *ServerResponse) SetBody(body string) *ServerResponse {
	sr.body = body
	return sr
}

func (sr *ServerResponse) SetStatus(status ResponseStatus) *ServerResponse {
	sr.status = status
	return sr
}

func (sr *ServerResponse) ToString() string {
	headers := ""
	for key, value := range sr.headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	response := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n%s", sr.status, headers, sr.body)
	return response
}

func handleRequest(conn net.Conn) {
	data := make([]byte, 1024)
	_, err := conn.Read(data)

	if err != nil {
		log.Fatalln("Error reading from connection", err)
	}

	parts := strings.Split(string(data), " ")
	request := ServerRequest{
		Method:  parts[0],
		Uri:     parts[1],
		Version: parts[2],
	}

	response := handlePath(request.Uri)
	fmt.Println("DEBUG: ", response.ToString())

	_, err = conn.Write([]byte(response.ToString()))
	if err != nil {
		log.Fatalln("Error while handling the request:", err)
	}
}

func handlePath(path string) ServerResponse {
	response := ServerResponse{headers: make(map[string]string)}
	pathParts := strings.Split(path, "/")

	if path == "/" {
		response.SetStatus(OK)
	} else if pathParts[1] == "echo" {
		response.
			SetHeader("Content-Type", "text/plain").
			SetHeader("Content-Length", strconv.Itoa(len(pathParts[2]))).
			SetStatus(OK).
			SetBody(pathParts[2])
	} else {
		response.SetStatus(NOT_FOUND)
	}

	return response
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

	handleRequest(connection)
}
