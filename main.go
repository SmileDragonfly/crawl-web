package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StoryInfo struct {
	ID          int
	Title       string
	Image       string
	Author      string
	Type        string
	Source      string
	Status      string
	Description string
	Contents    string
}

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
	// 1. Get 'the loai' links
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
	// 2. Loop get list story in each 'the loai'
	var storyInfo StoryInfo
	var storyLinks []string
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
		doc.Find("h3.truyen-title a").Each(func(i int, s *goquery.Selection) {
			// sHtml, _ := s.Html()
			// fmt.Println(sHtml)
			link, exists := s.Attr("href")
			if exists && !strings.HasPrefix(link, "javascript") {
				storyLinks = append(storyLinks, link)
				//fmt.Println(link)
			}
		})
	}
	// 3. Get contents of each story in storyLinks
	storyLink := storyLinks[0]
	storyRes, err := http.Get(storyLink)
	if err != nil {
		log.Fatal(err)
	}
	defer storyRes.Body.Close()
	doc, err = goquery.NewDocumentFromReader(storyRes.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Find title
	storyInfo.Title = doc.Find("h3.title").Text()
	// Find story image
	imgSrc, exists := doc.Find("div.book img").Attr("src")
	if exists {
		storyInfo.Image = imgSrc
	} else {
		fmt.Println("Không tìm thấy ảnh.")
	}
	// Find author
	authorName := doc.Find("a[itemprop='author']").Text()
	authorLink, exists := doc.Find("a[itemprop='author']").Attr("href")
	storyInfo.Author = authorName
	if exists {
		fmt.Println("Link tác giả:", authorLink)
	} else {
		fmt.Println("Không tìm thấy link.")
	}
	// Find type
	var sType string
	infoDiv := doc.Find("div.info")
	if infoDiv.Length() == 0 {
		fmt.Println("Không tìm thấy thẻ <div class='info'>")
		return
	}
	infoDiv.Find("a[itemprop='genre']").Each(func(i int, s *goquery.Selection) {
		fmt.Println("Genre:", s.Text())
		sType += s.Text()
	})
	storyInfo.Type = sType
	// Find source
	storyInfo.Source = infoDiv.Find("span.source").Text()
	fmt.Println("Nguồn:", storyInfo.Source)
}
