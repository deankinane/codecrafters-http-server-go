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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			return
		}

		HandleConnection(conn)
	}

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
		conn.Write([]byte(BuildResonseBody(200, echo)))
	} else if req.URL.Path == "/user-agent" {
		conn.Write([]byte(BuildResonseBody(200, req.UserAgent())))
	} else if req.URL.Path == "/" {
		conn.Write([]byte(BuildResonseBody(200, "")))
	} else {
		conn.Write([]byte(BuildResonseBody(404, "")))
	}
}

func BuildResonseBody(code int, body string) string {
	response := ResponseCode(code)

	if len(body) > 0 {
		response += "Content-Type: text/plain\r\n"
		response += "Content-Length: " + fmt.Sprint(len(body)) + "\r\n"
	}

	response += "\r\n"
	response += body
	return response
}

func ResponseCode(code int) string {
	codeStr := ""
	switch code {
	case 200:
		codeStr = "200 OK"
	case 404:
		codeStr = "404 Not Found"
	default:
		codeStr = "400 Bad Request"
	}

	return "HTTP/1.1 " + codeStr + "\r\n"
}
