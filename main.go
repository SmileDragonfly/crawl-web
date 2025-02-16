package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// fmt.Println("Hello")
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
	// Get 'the loai' links
	var links []string
	doc.Find("li.dropdown").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Thể loại") {
			// sHtml, _ := s.Html()
			// fmt.Println(sHtml)
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				link, exists := s.Attr("href")
				if exists && !strings.HasPrefix(link, "javascript") {
					links = append(links, link)
				}
			})
		}
	})
	// fmt.Print(links)
	// Loop get list story in each 'the loai'
	for _, link := range links {
		resLink, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		defer resLink.Body.Close()
		if resLink.StatusCode != 200 {
			log.Fatalf("Status code error: %d %s %s", res.StatusCode, res.Status, link)
		}
		doc, err := goquery.NewDocumentFromReader(resLink.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(doc.Html())
		var storyLinks []string
		doc.Find("h3.truyen-title a").Each(func(i int, s *goquery.Selection) {
			// sHtml, _ := s.Html()
			// fmt.Println(sHtml)
			link, exists := s.Attr("href")
			if exists && !strings.HasPrefix(link, "javascript") {
				links = append(storyLinks, link)
				fmt.Println(link)
			}
		})
	}
}
