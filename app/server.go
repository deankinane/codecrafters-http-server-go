package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

var param_directory string = ""

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	if len(os.Args) > 2 {
		param_directory = os.Args[2]
	}

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

		go HandleConnection(conn)
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
		conn.Write([]byte(BuildResponseText(200, echo)))
	} else if req.URL.Path == "/user-agent" {
		conn.Write([]byte(BuildResponseText(200, req.UserAgent())))
	} else if strings.HasPrefix(req.URL.Path, "/files/") {
		path := strings.Replace(req.URL.Path, "/files/", "", 1)
		conn.Write(([]byte(ServeFile(path))))
	} else if req.URL.Path == "/" {
		conn.Write([]byte(EmptyResponse(200)))
	} else {
		conn.Write([]byte(EmptyResponse(404)))
	}
}

func ServeFile(path string) string {
	content, err := os.ReadFile(param_directory + "/" + path)
	if err != nil {
		return EmptyResponse(404)
	}

	return BuildResponseBinary(200, content)
}

func BuildResponseBody(code int, content_type string, content_length int, body string) string {
	response := ResponseCode(code)

	if len(body) > 0 {
		response += "Content-Type: " + content_type + "\r\n"
		response += "Content-Length: " + fmt.Sprint(content_length) + "\r\n"
	}

	response += "\r\n"
	response += body
	return response
}

func BuildResponseText(code int, body string) string {
	return BuildResponseBody(code, "text/plain", len(body), body)
}

func BuildResponseBinary(code int, body []byte) string {
	return BuildResponseBody(code, "application/octet-stream", len(body), string(body))
}

func EmptyResponse(code int) string {
	return BuildResponseText(code, "")
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
