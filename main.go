package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Println("Hello")
	// Request the html page
	res, err := http.Get("https://truyenfull.vision")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("li.dropdown").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Thể loại") {

		}
	})
}
