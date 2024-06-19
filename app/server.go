package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"slices"
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

	if req.Method == "GET" {
		HandleGet(conn, req)
	} else if req.Method == "POST" {
		HandlePost(conn, req)
	}

}

func HandleGet(conn net.Conn, req *http.Request) {
	encoding := req.Header.Get("Accept-Encoding")
	if strings.HasPrefix(req.URL.Path, "/echo/") {
		fmt.Println(req.URL.Path)
		echo := strings.Replace(req.URL.Path, "/echo/", "", 1)
		conn.Write([]byte(BuildResponseText(200, echo, encoding)))
	} else if req.URL.Path == "/user-agent" {
		conn.Write([]byte(BuildResponseText(200, req.UserAgent(), encoding)))
	} else if strings.HasPrefix(req.URL.Path, "/files/") {
		path := strings.Replace(req.URL.Path, "/files/", "", 1)
		conn.Write([]byte(ServeFile(path, encoding)))
	} else if req.URL.Path == "/" {
		conn.Write([]byte(EmptyResponse(200)))
	} else {
		conn.Write([]byte(EmptyResponse(404)))
	}
}

func HandlePost(conn net.Conn, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/files/") {
		path := strings.Replace(req.URL.Path, "/files/", "", 1)
		body, err := io.ReadAll(req.Body)
		if err != nil {
			conn.Write([]byte(EmptyResponse(400)))
		}

		conn.Write([]byte(PostFile(path, body)))
	} else {
		conn.Write([]byte(EmptyResponse(404)))
	}
}

func ServeFile(path string, encoding string) string {
	content, err := os.ReadFile(param_directory + "/" + path)
	if err != nil {
		return EmptyResponse(404)
	}

	return BuildResponseBinary(200, content, encoding)
}

func BuildResponseBody(code int, content_type string, body []byte, encoding string) string {
	response := ResponseCode(code)

	encoding_types := strings.Fields(strings.ReplaceAll(encoding, ",", ""))
	content_encoding := ""
	if slices.Contains(encoding_types, "gzip") {
		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		w.Write(body)
		w.Close()
		body = b.Bytes()
		content_encoding = "Content-Encoding: gzip\r\n"
	}

	content_length := 0
	switch content_type {
	case "application/octet":
		content_length = len(body)
	default:
		content_length = len(string(body))
	}

	if len(body) > 0 {
		response += "Content-Type: " + content_type + "\r\n"
		response += "Content-Length: " + fmt.Sprint(content_length) + "\r\n"
		response += content_encoding
	}

	response += "\r\n"
	response += string(body)
	return response
}

func BuildResponseText(code int, body string, encoding string) string {
	return BuildResponseBody(code, "text/plain", []byte(body), encoding)
}

func BuildResponseBinary(code int, body []byte, encoding string) string {
	return BuildResponseBody(code, "application/octet-stream", body, encoding)
}

func EmptyResponse(code int) string {
	return BuildResponseText(code, "", "")
}

func ResponseCode(code int) string {
	codeStr := ""
	switch code {
	case 200:
		codeStr = "200 OK"
	case 201:
		codeStr = "201 Created"
	case 404:
		codeStr = "404 Not Found"
	default:
		codeStr = "400 Bad Request"
	}

	return "HTTP/1.1 " + codeStr + "\r\n"
}

func PostFile(path string, data []byte) string {
	err := os.WriteFile(param_directory+"/"+path, data, os.ModeAppend)
	if err != nil {
		return EmptyResponse(400)
	}

	return EmptyResponse(201)
}
