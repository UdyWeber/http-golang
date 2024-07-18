package http

import (
	"fmt"
)

type ServerResponse struct {
	status  ResponseStatus
	Headers map[string]string
	body    string
}

func (sr *ServerResponse) SetHeader(header string, value string) *ServerResponse {
	sr.Headers[header] = value
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
	for key, value := range sr.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	response := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n%s", sr.status, headers, sr.body)
	return response
}
