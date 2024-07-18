package http

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
)

func tryEncoding(request *ServerRequest, response *ServerResponse) *ServerResponse {
	defer fmt.Println("[INFO] Sending back response\n", response.ToString())

	requestEncoders, exists := request.Headers["Accept-Encoding"]

	if !exists {
		return response
	}

	validEncoders := getValidEncoders(requestEncoders)

	if len(validEncoders) == 0 {
		fmt.Println("[INFO] Possible encoding does not exist")
		return response
	}

	buffer := writeEncodedBuffer(validEncoders, response.body)

	fmt.Println("[INFO] Writing response body", buffer.String())
	return response.
		SetHeader("Content-Length", strconv.Itoa(len(buffer.String()))).
		SetHeader("Content-Encoding", validEncoders[0]).
		SetBody(buffer.String())
}

func writeEncodedBuffer(validEncoders []string, body string) *bytes.Buffer {
	fmt.Println("[INFO] Inserting encoding header", validEncoders[0])
	buffer := &bytes.Buffer{}
	writter := gzip.NewWriter(buffer)
	writter.Write([]byte(body))
	writter.Close()
	return buffer
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
