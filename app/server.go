package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type ResponseStatus string
type HttpMethod string

const (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

const (
	OK        ResponseStatus = "200 OK"
	NOT_FOUND ResponseStatus = "404 Not Found"
	CREATED   ResponseStatus = "201 Created"
)

var supportedEncodings = []string{"gzip", "deflate"}

var filesPath = ""

type ServerRequest struct {
	Method  HttpMethod
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

	body := sections[len(sections)-1]
	fmt.Println("[INFO] Request body:", body)
	value, exists := headers["Content-Length"]
	if exists {
		fmt.Println("[INFO] Reduced request body:", body, "with length:", value)
		length, _ := strconv.Atoi(value)
		body = body[:length]
	}

	return &ServerRequest{
		Method:  HttpMethod(requestLineParts[0]),
		Uri:     requestLineParts[1],
		Version: requestLineParts[2],
		Headers: headers,
		Body:    body,
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("[INFO] Handling the connection on thread ", os.Getpid())
	data := make([]byte, 1024)
	_, err := conn.Read(data)

	if err != nil {
		log.Println("[ERROR] Failed reading from connection", err)
		return
	}

	request := mountRequest(data)
	response := handleRequest(request)
	fmt.Println("[INFO] Sending back response\n", response.ToString())

	_, err = conn.Write([]byte(response.ToString()))
	if err != nil {
		log.Println("[ERROR] Failed while handling the request:", err)
		return
	}

	conn.Close()
}

func handleRequest(request *ServerRequest) *ServerResponse {
	response := &ServerResponse{headers: make(map[string]string)}
	pathParts := strings.Split(request.Uri, "/")

	fmt.Println("[DEBUG] Handling URI", request.Uri)
	fmt.Println("[DEBUG] Path Parts", pathParts)
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
	} else if pathParts[1] == "files" && request.Method == GET {
		fileName := filesPath + pathParts[2]
		fmt.Println("[INFO] Handling file", fileName)
		dat, err := os.ReadFile(fileName)

		if err != nil {
			return response.SetStatus(NOT_FOUND)
		}

		strData := string(dat)

		response.
			SetStatus(OK).
			SetHeader("Content-Type", "application/octet-stream").
			SetHeader("Content-Length", strconv.Itoa(len(strData))).
			SetBody(strData)
	} else if pathParts[1] == "files" && request.Method == POST {
		fileName := filesPath + pathParts[2]
		fmt.Println("[INFO] Creating file", fileName, "with content:", request.Body)
		err := os.WriteFile(fileName, []byte(request.Body), 0755)
		if err != nil {
			fmt.Println("[ERROR] Failed creating file", err)
		}
		response.SetStatus(CREATED)
	} else {
		response.SetStatus(NOT_FOUND)
	}

	return tryEncoding(request, response)
}

func tryEncoding(request *ServerRequest, response *ServerResponse) *ServerResponse {
	requestEncoders, exists := request.Headers["Accept-Encoding"]

	if !exists {
		return response
	}

	validEncoders := getValidEncoders(requestEncoders)

	if len(validEncoders) == 0 {
		fmt.Println("[INFO] Possible encoding does not exist")
		return response
	}

	fmt.Println("[INFO] Inserting encoding header", validEncoders[0])
	buffer := &bytes.Buffer{}
	writter := gzip.NewWriter(buffer)
	writter.Write([]byte(response.body))
	writter.Close()

	fmt.Println("[INFO] Writing response body", buffer.String())

	return response.
		SetHeader("Content-Length", strconv.Itoa(len(buffer.String()))).
		SetHeader("Content-Encoding", validEncoders[0]).
		SetBody(buffer.String())
}

func getValidEncoders(value string) []string {
	var validEncoders []string
	possibleEncodings := strings.Split(value, ",")
	fmt.Println("[INFO] Handling Accept-Encoding ", possibleEncodings)

	for _, supportedEncoding := range supportedEncodings {
		for _, possibleEncoding := range possibleEncodings {
			fmt.Println("[INFO] Comparing", possibleEncoding, " with", supportedEncoding)
			if strings.Compare(strings.TrimSpace(possibleEncoding), supportedEncoding) == 0 {
				fmt.Println("[INFO] Possible encoding exists", supportedEncoding)
				validEncoders = append(validEncoders, supportedEncoding)
			}
		}
	}
	return validEncoders
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	fmt.Println("[DEBUG] Program params", os.Args[1:])

	if len(os.Args) > 2 && os.Args[1] == "--directory" && os.Args[2] != "" {
		fmt.Println("[INFO] Setting files directory to ", os.Args[2])
		filesPath = os.Args[2]
	}

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
