package protocols

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type HTTP struct {
	host       string
	path       string
	port       string
	protocol   string
	statusCode string

	RespHeaders map[string]string
	ReqHeaders  map[string]string
}

func NewHTTP(host string, path string, port string, protocol string) *HTTP {
	return &HTTP{
		host:        host,
		path:        path,
		port:        port,
		protocol:    protocol,
		statusCode:  "",
		RespHeaders: make(map[string]string),
		ReqHeaders:  make(map[string]string),
	}
}

func (h *HTTP) Host() string {
	return h.host
}

func (h *HTTP) Path() string {
	return h.path
}

func (h *HTTP) SetHeader(key string, value string) {
	h.ReqHeaders[key] = value
}

func (h *HTTP) Request() (string, error) {
	var err error
	var conn net.Conn
	if h.protocol == "http" {
		conn, err = net.Dial("tcp", h.host+":"+h.port)
		if err != nil {
			return "", err
		}
		defer conn.Close()
	} else if h.protocol == "https" {
		conn, err = tls.Dial("tcp", h.host+":"+h.port, &tls.Config{
			// FIXME: Bad for production
			InsecureSkipVerify: true,
		})
		if err != nil {
			return "", err
		}
		defer conn.Close()
	} else if h.protocol == "file" {
		f, err := os.ReadFile(h.path)
		if err != nil {
			return "", err
		}
		return string(f), nil
	} else if h.protocol == "data" {
		return h.path, nil
	}
	reqHeaders := ""
	for k, v := range h.ReqHeaders {
		reqHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	if _, exist := h.ReqHeaders["connection"]; !exist {
		reqHeaders += "Connection: close\r\n"
	}
	reqHeaders += fmt.Sprintf("Host: %s\r\n", h.host)
	reqHeaders += fmt.Sprintf("User-Agent: %s\r\n", "Ignis/SimpleBrowser")
	fmt.Fprintf(conn, "GET %s HTTP/1.1\r\n%s\r\n", h.path, reqHeaders)

	responses := make([]string, 0)
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		responses = append(responses, message)
	}

	statusLine := strings.SplitN(responses[0], " ", 3)
	if len(statusLine) != 3 {
		return "", errors.New("something really bad happened")
	}
	version, statusCode, explanation := statusLine[0], statusLine[1], statusLine[2]
	// TODO: Put in response struct or something later
	fmt.Println(version, statusCode, explanation)
	h.statusCode = statusCode

	i := 1
	for i < len(responses) {
		line := responses[i]
		i++
		if line == "\r\n" {
			break
		}
		keyVal := strings.SplitN(line, ":", 2)
		key, value := keyVal[0], keyVal[len(keyVal)-1]
		h.RespHeaders[strings.ToLower(key)] = strings.TrimSpace(value)
	}

	if _, exist := h.RespHeaders["transfer-encoding"]; exist {
		return "", errors.ErrUnsupported
	}
	if _, exist := h.RespHeaders["content-encoding"]; exist {
		return "", errors.ErrUnsupported
	}

	body := responses[i:]

	return strings.Join(body, ""), nil
}

func (h *HTTP) StatusCode() string {
	return h.statusCode
}

func (h *HTTP) ResponseHeaders() map[string]string {
	return h.RespHeaders
}

func (h *HTTP) Protocol() string {
	return h.protocol
}

func (h *HTTP) RequestHeaders() map[string]string {
	return h.ReqHeaders
}

func (h *HTTP) SetHeaders(headers map[string]string) {
	h.ReqHeaders = headers
}
