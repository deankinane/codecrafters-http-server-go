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
		response := "HTTP/1.1 200 OK\r\n"
		response += "Content-Type: text/plain\r\n"
		response += "Content-Length: " + fmt.Sprint(len(echo)) + "\r\n"
		response += "\r\n"
		response += echo
		conn.Write([]byte(response))
	}

}
