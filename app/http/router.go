package http

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/globals"
	"os"
	"strconv"
	"strings"
)

func HandleRequest(request *ServerRequest) *ServerResponse {
	response := &ServerResponse{Headers: make(map[string]string)}
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
		fileName := globals.FilesPath + pathParts[2]
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
		fileName := globals.FilesPath + pathParts[2]
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
