package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
)

// 17:06:43 2025-03-30: Unused function
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
		log.Println(err)
		return
	}
	links := strings.Split(string(data), "\r\n")
	// Open DB
	db, err := OpenDB("127.0.0.1", 3306, "truyenfull", "root", "123@123A")
	if err != nil {
		log.Println(err)
		return
	}
	// [Test] 100 link
	links = links[:100]
	for i, link := range links {
		log.Printf("[%d]Get story %s\n", i, link)
		GetStory(link, db)
	}
}

func GetStory(link string, db *gorm.DB) {
	res, err := http.Get(link)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	// 22:28:01 2025-03-31: Check if story exist or not
	storyInfo := GetStoryByUrl(db, link)
	if storyInfo.Url == "" {
		storyInfo.Url = link
		// Find title
		storyInfo.Title = doc.Find("h3.title").Text()
		// Find story image
		imgSrc, exists := doc.Find("div.book img").Attr("src")
		if exists {
			storyInfo.Image = imgSrc
		} else {
			log.Printf("[%s]Cannot find image", link)
			return
		}
		// Find author
		authorName := doc.Find("a[itemprop='author']").Text()
		authorLink, _ := doc.Find("a[itemprop='author']").Attr("href")
		storyInfo.Author = authorName
		storyInfo.AuthorUrl = authorLink
		// Find type
		var sType []string
		infoDiv := doc.Find("div.info")
		if infoDiv.Length() > 0 {
			infoDiv.Find("a[itemprop='genre']").Each(func(i int, s *goquery.Selection) {
				//fmt.Println("Genre:", s.Text())
				sType = append(sType, s.Text())
			})
			// Chuyển danh sách sang JSON string
			genresJSON, _ := json.Marshal(sType)
			storyInfo.Type = string(genresJSON)
			// Find source
			storyInfo.Source = infoDiv.Find("span.source").Text()
			//fmt.Println("Nguồn:", storyInfo.Source)
			// Find status
			status := infoDiv.Find("span.text-success, span.text-primary").Text()
			storyInfo.Status = status
		}
		//fmt.Println("Trạng thái:", storyInfo.Status)
		// Find rating
		rate, exist := doc.Find("div.rate-holder").Attr("data-score")
		if exist {
			storyInfo.Rate = rate
		}
		// Find rating count
		ratingCount := doc.Find("span[itemprop='ratingCount']").Text()
		storyInfo.RatingCount = ratingCount
		//fmt.Println("Rating:", storyInfo.Rate, "-RatingCount:", storyInfo.RatingCount)
		// Find description
		description, err := doc.Find("div.desc-text").Html()
		if err != nil {
			log.Printf("[%s]Cannot file description: %s\n", link, err.Error())
			return
		}
		storyInfo.Description = description
		// Save info to DB
		err = SaveStory(db, &storyInfo)
		if err != nil {
			log.Printf("[%s]Cannot save story info to DB: %s\n", link, err.Error())
			log.Println(err)
			return
		}
	} else {
		log.Printf("Url existed: %s => Update chapters\n", storyInfo.Url)
	}
	// Get all chapters of this story
	GetChapter(doc, storyInfo.ID, db)
}

func GetChapter(doc *goquery.Document, storyID int, db *gorm.DB) {
	log.Println("Begin get chapters")
	// Find list page
	var lastPage int = 0
	var lastPageLink string
	found := doc.Find("ul.pagination.pagination-sm li a").Each(func(index int, item *goquery.Selection) {
		link, exists := item.Attr("href")
		text := item.Text()
		// Dùng regex để lấy số trang từ URL
		re := regexp.MustCompile(`trang-(\d+)`)
		match := re.FindStringSubmatch(link)
		var currPage int
		if len(match) > 1 {
			currPage, _ = strconv.Atoi(match[1])
		}
		// Nếu là link trang cuối (có chữ "Cuối" hoặc số lớn nhất)
		if exists && strings.Contains(text, "Cuối") {
			lastPageLink = link
			lastPage = currPage
		} else {
			if currPage > lastPage {
				lastPage = currPage
				lastPageLink = link
			}
		}
	})
	// Find list chapter
	getChapterContent := func(doc *goquery.Document, storyID int) {
		doc.Find(".list-chapter li a").Each(func(index int, element *goquery.Selection) {
			// Lấy link chương
			link, exists := element.Attr("href")
			if !exists {
				return
			}
			// Lấy tên chương từ `title`
			title, _ := element.Attr("title")
			// Lấy chapter number từ link
			re := regexp.MustCompile(`/chuong-(\d+)/`)
			matches := re.FindStringSubmatch(link)
			chapterNumber := 0
			if len(matches) > 1 {
				if num, err := strconv.Atoi(matches[1]); err == nil {
					chapterNumber = num
				}
			}
			// 23:32:55 2025-03-31: Check chapter existed in DB or not
			chapter := GetChapterByUrl(db, link)
			if chapter.Url != "" {
				log.Printf("[%s]Chapter %d existed\n", link, chapterNumber)
				return
			}
			log.Printf("[%s]Crawling chapter %d\n", link, chapterNumber)
			// Get content
			res, err := http.Get(link)
			if err != nil {
				log.Printf("[%s]Cannot http get chapter: %s\n", link, err.Error())
				return
			}
			defer res.Body.Close()
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Printf("[%s]Cannot create document: %s\n", link, err.Error())
				return
			}
			html, err := doc.Find("#chapter-c").Html()
			if err != nil {
				log.Printf("[%s]Cannot find chapter content: %s\n", link, err.Error())
				return
			}
			chapter = Chapter{
				StoryID:       storyID,
				ChapterNumber: chapterNumber,
				Title:         title,
				Url:           link,
				Content:       html,
			}
			err = SaveChapter(db, chapter)
			if err != nil {
				log.Printf("[%s]Cannot save chapter: %s\n", link, err.Error())
				return
			}
		})
	}
	if found.Size() > 0 {
		// Duyet qua tung page
		baseURL := strings.Replace(lastPageLink, fmt.Sprintf("trang-%d", lastPage), "trang-%d", 1)
		for i := range lastPage {
			pageURL := fmt.Sprintf(baseURL, i+1)
			pageRes, err := http.Get(pageURL)
			if err != nil {
				log.Printf("[%s]Cannot http get page url: %s\n", pageURL, err.Error())
				return
			}
			defer pageRes.Body.Close()
			pageDoc, err := goquery.NewDocumentFromReader(pageRes.Body)
			if err != nil {
				log.Printf("[%s]Cannot create new doc from page url: %s\n", pageURL, err.Error())
				return
			}
			// Duyệt qua từng chương
			getChapterContent(pageDoc, storyID)
		}
	} else {
		getChapterContent(doc, storyID)
	}
}
