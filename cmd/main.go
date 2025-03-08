package main

import (
	"fmt"
	"log"

	"github.com/Lunarisnia/simple-browser/internal/url"
)

func main() {
	u, err := url.New("https://browser.engineering/examples/example1-simple.html")
	if err != nil {
		log.Fatal(err)
	}
	content, err := u.Request()
	if err != nil {
		log.Fatal(err)
	}
	// Show(content)
	body := Show(content)
	fmt.Println(body)
	// fmt.Println("Content: ", content)
}

func Show(body string) string {
	parsedBody := ""
	inTag := false
	for _, c := range body {
		if c == '<' {
			inTag = true
		} else if c == '>' {
			inTag = false
		} else if !inTag {
			// fmt.Print(string(c))
			parsedBody += string(c)
		}
	}

	return parsedBody
}
