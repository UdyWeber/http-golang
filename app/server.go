package main

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/globals"
	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"log"
	"net"
	"os"
)

func handleExecutionArgs() {
	fmt.Println("[DEBUG] Program params", os.Args[1:])

	if len(os.Args) > 2 && os.Args[1] == "--directory" && os.Args[2] != "" {
		fmt.Println("[INFO] Setting files directory to ", os.Args[2])
		globals.FilesPath = os.Args[2]
	}
}

func startServingConnections() {
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

func handleConnection(conn net.Conn) {
	fmt.Println("[INFO] Handling the connection on thread ", os.Getpid())
	data := make([]byte, 1024)
	_, err := conn.Read(data)

	if err != nil {
		log.Println("[ERROR] Failed reading from connection", err)
		return
	}

	request := http.MountRequest(data)
	response := http.HandleRequest(request)

	_, err = conn.Write([]byte(response.ToString()))
	if err != nil {
		log.Println("[ERROR] Failed while handling the request:", err)
		return
	}

	conn.Close()
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	handleExecutionArgs()
	startServingConnections()
}
