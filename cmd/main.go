package main

import (
	"fmt"
	"log"

	"github.com/Lunarisnia/simple-browser/internal/url"
)

func main() {
	// u, err := url.New("data:text/html,&lt;div&gt;;;;")
	// u, err := url.New("http://browser.engineering/http.html")
	u, err := url.New("https://example.org/index.html")
	u.SetHeader("Accept-Encoding", "gzip")
	if err != nil {
		log.Fatal(err)
	}
	body, err := url.Load(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(body)
	// fmt.Println("Content: ", content)
}
