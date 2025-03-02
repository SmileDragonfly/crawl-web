package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetStoryByChapterNumber(sourceUrl string) {
	// 1. Get source html
	res, err := http.Get(sourceUrl)
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
	// 2. Get 'Phân loại theo Chương' links
	fileByChapter, err := os.Create("LinkByChapter.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fileByChapter.Close()
	var links []string
	doc.Find("li.dropdown").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Phân loại theo Chương") {
			// sHtml, _ := s.Html()
			// fmt.Println(sHtml)
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				link, exists := s.Attr("href")
				if exists && !strings.HasPrefix(link, "javascript") {
					links = append(links, link)
					fileByChapter.WriteString(link)
					fileByChapter.WriteString("\n")
				}
			})
		}
	})
	// fmt.Println(links)
	// 3. Get story in each link
	fileStoryLink, err := os.Create("StoryLink.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fileStoryLink.Close()
	for _, link := range links {
		resLink, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		defer resLink.Body.Close()
		docLink, err := goquery.NewDocumentFromReader(resLink.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Find all pagination pages
		lastPage, lastPageLink := func(doc *goquery.Document) (int, string) {
			var lastPage int
			var lastPageLink string
			doc.Find("li a").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				if strings.Contains(text, "Cuối") {
					lastPageLink, _ = s.Attr("href")
					// Find last page index
					re := regexp.MustCompile(`trang-(\d+)/`)
					match := re.FindStringSubmatch(lastPageLink)
					if len(match) > 1 {
						lastPage, _ = strconv.Atoi(match[1])
					}
				}
			})
			return lastPage, lastPageLink
		}(docLink)
		fmt.Println("lastPageLink:", lastPageLink)
		fmt.Println("lastPage:", lastPage)
		// Scan each page to get all story links
		storyList := func(lastPage int, lastPageLink string) []string {
			var storyList []string
			for i := range lastPage {
				pageUrl := strings.Replace(lastPageLink, strconv.Itoa(lastPage), strconv.Itoa(i), 1)
				resPage, err := http.Get(pageUrl)
				if err != nil {
					log.Fatal(err)
				}
				defer resPage.Body.Close()
				docPage, err := goquery.NewDocumentFromReader(resPage.Body)
				if err != nil {
					log.Fatal(err)
				}
				docPage.Find("h3.truyen-title a").Each(func(i int, s *goquery.Selection) {
					storyUrl, _ := s.Attr("href")
					storyList = append(storyList, storyUrl)
					fileStoryLink.WriteString(storyUrl)
					fileStoryLink.WriteString("\n")
				})
				time.Sleep(time.Duration(100))
			}
			return storyList
		}(lastPage, lastPageLink)
		fmt.Println(storyList)
	}
}

func GetStoryInfoFromFile(filePath string) {
	// Read all lines from file then store to string array
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	links := strings.Split(string(data), "\n")
	// Open DB
	db, err := OpenDB("127.0.0.1", 5432, "truyenfull", "postgres", "123@123A")
	if err != nil {
		log.Fatal(err)
	}
	// [Test] 100 link
	links = links[:100]
	for _, link := range links {
		res, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Get story info then save to postgresql
		func(doc *goquery.Document) {
			var storyInfo StoryInfo
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
			authorLink, _ := doc.Find("a[itemprop='author']").Attr("href")
			storyInfo.Author = authorName
			storyInfo.AuthorUrl = authorLink
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
			// Find status
			status := infoDiv.Find("span.text-success, span.text-primary").Text()
			storyInfo.Status = status
			fmt.Println("Trạng thái:", storyInfo.Status)
			// Find rating
			rate, exist := doc.Find("div.rate-holder").Attr("data-score")
			if exist {
				storyInfo.Rate = rate
			}
			// Find rating count
			ratingCount := doc.Find("span[itemprop='ratingCount']").Text()
			storyInfo.RatingCount = ratingCount
			fmt.Println("Rating:", storyInfo.Rate, "-RatingCount:", storyInfo.RatingCount)
			// Find description
			description, err := doc.Find("div.desc-text").Html()
			if err != nil {
				log.Fatal(err)
			}
			storyInfo.Description = description
			// Save info to DB
			err = SaveStory(db, storyInfo)
			if err != nil {
				log.Println(err)
			}
		}(doc)
	}
}
