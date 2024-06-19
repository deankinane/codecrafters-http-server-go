package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		return
	}

	HandleConnection(conn)
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error readind request. ", err.Error())
	}

	if strings.HasPrefix(req.URL.Path, "/echo/") {
		fmt.Println(req.URL.Path)
		echo := strings.Replace(req.URL.Path, "/echo/", "", 1)
		conn.Write([]byte(BuildResonseBody(echo)))
	} else if req.URL.Path == "/user-agent" {
		conn.Write([]byte(BuildResonseBody(req.UserAgent())))
	} else if req.URL.Path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func BuildResonseBody(body string) string {
	response := "HTTP/1.1 200 OK\r\n"
	response += "Content-Type: text/plain\r\n"
	response += "Content-Length: " + fmt.Sprint(len(body)) + "\r\n"
	response += "\r\n"
	response += body
	return response
}
