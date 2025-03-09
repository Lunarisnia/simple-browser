package main

import (
	"fmt"
	"log"

	"github.com/Lunarisnia/simple-browser/internal/url"
)

func main() {
	// u, err := url.New("data:text/html,&lt;div&gt;;;;")
	u, err := url.New("view-source:https://example.org")
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
