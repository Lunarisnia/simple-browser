package url

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

type u struct {
	host     string
	path     string
	port     string
	protocol string
}

type URL interface {
	Request() (string, error)
	Host() string
	Path() string
}

func New(rawURL string) (URL, error) {
	parsedURL, err := parseRawURL(rawURL)
	if err != nil {
		return nil, err
	}
	return parsedURL, nil
}

func parseRawURL(rawURL string) (URL, error) {
	scheme := strings.SplitN(rawURL, "://", 2)
	protocol := scheme[0]
	if protocol != "http" && protocol != "https" {
		return nil, errors.New("invalid protocol")
	}

	scheme = strings.SplitN(scheme[len(scheme)-1], "/", 2)
	if len(scheme) == 1 {
		scheme = append(scheme, "")
	}
	path := "/" + scheme[len(scheme)-1]

	host := scheme[0]
	port := ""
	if strings.Contains(scheme[0], ":") {
		scheme = strings.SplitN(scheme[0], ":", 2)
		host = scheme[0]
		port = scheme[len(scheme)-1]
	}

	if port == "" && protocol == "http" {
		port = "80"
	} else if port == "" && protocol == "https" {
		port = "443"
	}

	parsedURL := u{
		host:     host,
		path:     path,
		port:     port,
		protocol: protocol,
	}
	fmt.Println(parsedURL.port, parsedURL.host, parsedURL.path)

	return &parsedURL, nil
}

func (ur *u) Host() string {
	return ur.host
}

func (ur *u) Path() string {
	return ur.path
}

func (ur *u) Request() (string, error) {
	var err error
	var conn net.Conn
	if ur.protocol == "http" {
		conn, err = net.Dial("tcp", ur.host+":"+ur.port)
		if err != nil {
			return "", err
		}
		defer conn.Close()
	} else {
		conn, err = tls.Dial("tcp", ur.host+":"+ur.port, &tls.Config{
			// FIXME: Bad for production
			InsecureSkipVerify: true,
		})
		if err != nil {
			return "", err
		}
		defer conn.Close()
	}
	fmt.Fprintf(conn, "GET %s HTTP/1.0\r\nHost: %s\r\n\r\n", ur.path, ur.host)

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
	fmt.Println(version, statusCode, explanation)

	i := 1
	headers := make(map[string]string)
	for i < len(responses) {
		line := responses[i]
		i++
		if line == "\r\n" {
			break
		}
		keyVal := strings.SplitN(line, ":", 2)
		key, value := keyVal[0], keyVal[len(keyVal)-1]
		headers[strings.ToLower(key)] = strings.TrimSpace(value)
	}

	if _, exist := headers["transfer-encoding"]; exist {
		return "", errors.ErrUnsupported
	}
	if _, exist := headers["content-encoding"]; exist {
		return "", errors.ErrUnsupported
	}

	body := responses[i:]

	return strings.Join(body, ""), nil
}
