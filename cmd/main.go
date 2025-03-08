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
	body, err := url.Load(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(body)
	// fmt.Println("Content: ", content)
}
