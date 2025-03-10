package url

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Lunarisnia/simple-browser/internal/caches"
	"github.com/Lunarisnia/simple-browser/internal/url/protocols"
)

var whitelistedProtocol map[string]bool = map[string]bool{
	"http":        true,
	"https":       true,
	"file":        true,
	"data":        true,
	"view-source": true,
}

var escapeCharMap map[string]rune = map[string]rune{
	"lt": '<',
	"gt": '>',
}

var cache caches.CacheBox

func init() {
	cache = caches.New()
}

type URL interface {
	Request() (string, error)
	Host() string
	Path() string
	Protocol() string
	StatusCode() string
	ResponseHeaders() map[string]string
	SetHeader(key string, value string)
	RequestHeaders() map[string]string
	SetHeaders(headers map[string]string)
}

func New(rawURL string) (URL, error) {
	parsedURL, err := parseRawURL(rawURL)
	if err != nil {
		return nil, err
	}
	return parsedURL, nil
}

func parseRawURL(rawURL string) (URL, error) {
	var clientProtocol string
	var scheme []string
	if strings.HasPrefix(rawURL, "data:") {
		scheme = strings.SplitN(rawURL, ":", 2)
	} else if strings.HasPrefix(rawURL, "view-source:") {
		scheme = strings.SplitN(rawURL, ":", 2)
		clientProtocol = scheme[0]
		scheme = strings.SplitN(scheme[len(scheme)-1], "://", 2)
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

	if clientProtocol == "view-source" {
		parsedURL := protocols.NewViewSource(host, path, port, protocol)
		return parsedURL, nil
	}

	parsedURL := protocols.NewHTTP(host, path, port, protocol)
	return parsedURL, nil
}

func Lex(body string) string {
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
	redirectionLimit := 10
	redirectionCounter := 0
	for redirectionCounter < redirectionLimit {
		content, err := u.Request()
		if err != nil {
			return "", err
		}
		if u.StatusCode() == "301" {
			newLocation := u.ResponseHeaders()["location"]
			if newLocation[0] == '/' {
				newLocation = u.Protocol() + "://" + u.Host() + newLocation
			}
			oldHeaders := u.RequestHeaders()
			u, err = New(newLocation)
			u.SetHeaders(oldHeaders)
			if err != nil {
				return "", err
			}
			redirectionCounter++
		} else {
			if v, exist := u.ResponseHeaders()["cache-control"]; exist {
				if v != "no-store" && strings.Contains(v, "max-age") {
					maxAgeStr := strings.TrimPrefix(v, "max-age=")
					maxAge, err := strconv.Atoi(maxAgeStr)
					if err != nil {
						return "", err
					}
					cachePath := map[string]string{
						u.Path(): content,
					}
					cache.Set(u.Host(), cachePath, maxAge)
				}
			}
			return Lex(content), nil
		}

	}
	return "", errors.New("too much redirects")
}
