package protocols

import "strings"

var charToEscape map[rune]string = map[rune]string{
	'<': "&lt;",
	'>': "&gt;",
}

type ViewSource struct {
	host     string
	path     string
	port     string
	protocol string
}

func NewViewSource(host string, path string, port string, protocol string) *ViewSource {
	return &ViewSource{
		host:     host,
		path:     path,
		port:     port,
		protocol: protocol,
	}
}

func (v *ViewSource) Request() (string, error) {
	httpProtocol := NewHTTP(v.host, v.path, v.port, v.protocol)
	body, err := httpProtocol.Request()
	if err != nil {
		return "", err
	}
	return escapeCharacters(body), nil
}

func (v *ViewSource) Host() string {
	return v.host
}

func (v *ViewSource) Path() string {
	return v.path
}

func escapeCharacters(body string) string {
	body = strings.ReplaceAll(body, "<", "&lt;")
	body = strings.ReplaceAll(body, ">", "&gt;")
	return body
}
