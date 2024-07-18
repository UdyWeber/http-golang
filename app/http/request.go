package http

import (
	"fmt"
	"strconv"
	"strings"
)

type ServerRequest struct {
	Method  HttpMethod
	Uri     string
	Version string
	Body    string
	Headers map[string]string
}

func MountRequest(data []byte) *ServerRequest {
	fullRequestString := string(data)
	sections := strings.Split(fullRequestString, "\r\n")

	requestLineParts := strings.Split(sections[0], " ")

	headers := mountHeaders(sections)

	body := formatBody(sections, headers)

	return &ServerRequest{
		Method:  HttpMethod(requestLineParts[0]),
		Uri:     requestLineParts[1],
		Version: requestLineParts[2],
		Headers: headers,
		Body:    body,
	}
}

func formatBody(sections []string, headers map[string]string) string {
	body := sections[len(sections)-1]
	fmt.Println("[INFO] Request body:", body)
	value, exists := headers["Content-Length"]
	if exists {
		fmt.Println("[INFO] Reduced request body:", body, "with length:", value)
		length, _ := strconv.Atoi(value)
		body = body[:length]
	}
	return body
}

func mountHeaders(sections []string) map[string]string {
	headers := make(map[string]string)

	for _, header := range sections[1 : len(sections)-2] {
		headerParts := strings.Split(header, ": ")
		headers[headerParts[0]] = headerParts[1]
	}

	return headers
}
