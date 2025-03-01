package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func OpenDB(dbIp string, dbPort int, dbName string, dbUser string, dbPassword string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbIp, dbPort, dbUser, dbPassword, dbName)
	return sql.Open("postgres", dsn)
}

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
	var links []string
	doc.Find("li.dropdown").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Phân loại theo Chương") {
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
	// fmt.Println(links)
	// 3. Get story in each link
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
				})
				// time.Sleep(time.Duration(100))
			}
			return storyList
		}(lastPage, lastPageLink)
		fmt.Println(storyList)
	}
}
