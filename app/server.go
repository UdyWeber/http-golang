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
	Body    string
	Headers map[string]string
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

func mountRequest(data []byte) *ServerRequest {
	fullRequestString := string(data)

	sections := strings.Split(fullRequestString, "\r\n")
	requestLineParts := strings.Split(sections[0], " ")

	headers := make(map[string]string)

	for _, header := range sections[1 : len(sections)-2] {
		headerParts := strings.Split(header, ": ")

		headers[headerParts[0]] = headerParts[1]
	}

	return &ServerRequest{
		Method:  requestLineParts[0],
		Uri:     requestLineParts[1],
		Version: requestLineParts[2],
		Headers: headers,
		Body:    sections[len(sections)-1],
	}
}

func handleConnection(conn net.Conn) {
	data := make([]byte, 1024)
	_, err := conn.Read(data)

	if err != nil {
		log.Println("[ERROR] Failed reading from connection", err)
		return
	}

	request := mountRequest(data)
	response := handleRequest(request)
	fmt.Println("[DEBUG] ", response.ToString())

	_, err = conn.Write([]byte(response.ToString()))
	if err != nil {
		log.Println("[ERROR] Failed while handling the request:", err)
		return
	}
}

func handleRequest(request *ServerRequest) *ServerResponse {
	response := &ServerResponse{headers: make(map[string]string)}
	pathParts := strings.Split(request.Uri, "/")

	fmt.Println("DEBUG: ", request.Uri)
	if request.Uri == "/" {
		response.SetStatus(OK)
	} else if request.Uri == "/user-agent" {
		userAgent := request.Headers["User-Agent"]
		response.
			SetHeader("Content-Type", "text/plain").
			SetHeader("Content-Length", strconv.Itoa(len(userAgent))).
			SetStatus(OK).
			SetBody(userAgent)
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

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := l.Accept()

		fmt.Println("[DEBUG] Connection accepted: ", connection.RemoteAddr().String())
		if err != nil {
			fmt.Println("[ERROR] Failed accepting connection: ", err.Error())
			continue
		}

		go handleConnection(connection)
	}
}
