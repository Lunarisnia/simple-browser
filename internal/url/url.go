package url

import (
	"errors"
	"strings"

	"github.com/Lunarisnia/simple-browser/internal/url/protocols"
)

var whitelistedProtocol map[string]bool = map[string]bool{
	"http":  true,
	"https": true,
	"file":  true,
	"data":  true,
}

var escapeCharMap map[string]rune = map[string]rune{
	"lt": '<',
	"gt": '>',
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
	var scheme []string
	if strings.HasPrefix(rawURL, "data:") {
		scheme = strings.SplitN(rawURL, ":", 2)
	} else {
		scheme = strings.SplitN(rawURL, "://", 2)
	}
	protocol := scheme[0]
	if _, exist := whitelistedProtocol[protocol]; !exist {
		return nil, errors.New("invalid protocol")
	}

	if protocol == "file" {
		parsedURL := protocols.NewHTTP("", scheme[len(scheme)-1], "", protocol)
		return parsedURL, nil
	}

	if protocol == "data" {
		scheme = strings.SplitN(rawURL, ",", 2)
		parsedURL := protocols.NewHTTP(scheme[0], scheme[len(scheme)-1], "", protocol)
		return parsedURL, nil
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

	parsedURL := protocols.NewHTTP(host, path, port, protocol)
	return parsedURL, nil
}

func Show(body string) string {
	parsedBody := ""
	escapeChar := ""
	inTag := false
	isEscaped := false
	for _, c := range body {
		if c == '<' {
			inTag = true
		} else if c == '>' {
			inTag = false
		} else if c == '&' {
			isEscaped = true
		} else if c == ';' && isEscaped {
			isEscaped = false
		} else if isEscaped {
			escapeChar += string(c)
			if v, exist := escapeCharMap[escapeChar]; exist {
				parsedBody += string(v)
				escapeChar = ""
			}
		} else if !inTag && !isEscaped {
			parsedBody += string(c)
		}
	}

	return parsedBody
}

func Load(u URL) (string, error) {
	content, err := u.Request()
	if err != nil {
		return "", err
	}

	return Show(content), nil
}
