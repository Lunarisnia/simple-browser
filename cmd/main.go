package main

import (
	"log"

	"github.com/Lunarisnia/simple-browser/internal/browser"
	"github.com/Lunarisnia/simple-browser/internal/url"
)

func main() {
	// // u, err := url.New("data:text/html,&lt;div&gt;;;;")
	// // u, err := url.New("http://browser.engineering/http.html")
	// u, err := url.New("https://example.org/index.html")
	// u.SetHeader("Accept-Encoding", "chunked")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// body, err := url.Load(u)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(body)
	// // fmt.Println("Content: ", content)

	ignis := browser.New(800, 600)
	// u, err := url.New("https://browser.engineering/examples/xiyouji.html")
	u, err := url.New("data:text/html,hello there")
	if err != nil {
		log.Fatal(err)
	}
	ignis.Load(u)
	ignis.Run()
}
