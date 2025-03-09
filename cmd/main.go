package main

import (
	"log"

	"github.com/Lunarisnia/simple-browser/internal/url"
)

func main() {
	// u, err := url.New("data:text/html,&lt;div&gt;;;;")
	u, err := url.New("http://browser.engineering/redirect3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = url.Load(u)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(body)
	// fmt.Println("Content: ", content)
}
