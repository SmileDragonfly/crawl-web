package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StoryInfo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Author      string `json:"author"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	Status      string `json:"status"`
	Rate        string `json:"rate"`
	RatingCount string `json:"rating_count"`
	Description string `json:"description"`
	Contents    string `json:"contents"`
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
	/* Comment for fast test a story link
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
	*/
	// 3. Get contents of each story in storyLinks
	// storyLink := "https://truyenfull.vision/ba-chu-thien-ha-phong-than/"
	storyLink := "https://truyenfull.vision/chang-chat-phac-nang-xau-xi/"
	storyRes, err := http.Get(storyLink)
	if err != nil {
		log.Fatal(err)
	}
	defer storyRes.Body.Close()
	doc, err = goquery.NewDocumentFromReader(storyRes.Body)
	if err != nil {
		log.Fatal(err)
	}
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
	description = strings.ReplaceAll(description, "<br>", "\n")
	description = strings.ReplaceAll(description, "<br/>", "\n") // Handle self-closing <br>
	description = strings.ReplaceAll(description, "&nbsp;", " ") // Remove any remaining HTML entities (like &nbsp;)(nbsp: non breaking space)
	description = strings.ReplaceAll(description, "&#34;", `"`)
	storyInfo.Description = description
	// Find list page
	// Tìm trang cuối cùng (có chữ "Cuối")
	var lastPage int
	var lastPageLink string
	found := doc.Find("ul.pagination.pagination-sm li a").Each(func(index int, item *goquery.Selection) {
		link, exists := item.Attr("href")
		text := item.Text()

		// Nếu là link trang cuối (có chữ "Cuối" hoặc số lớn nhất)
		if exists && strings.Contains(text, "Cuối") {
			lastPageLink = link
			// Dùng regex để lấy số trang từ URL
			re := regexp.MustCompile(`trang-(\d+)`)
			match := re.FindStringSubmatch(link)
			if len(match) > 1 {
				lastPage, _ = strconv.Atoi(match[1])
			}
		}
	})

	// Find list chapter
	var chapters []map[string]string
	if found.Size() > 0 {
		// Duyet qua tung page
		baseURL := strings.Replace(lastPageLink, fmt.Sprintf("trang-%d", lastPage), "trang-%d", 1)
		for i := range lastPage {
			pageURL := fmt.Sprintf(baseURL, i+1)
			pageRes, err := http.Get(pageURL)
			if err != nil {
				log.Fatal(err)
			}
			defer pageRes.Body.Close()
			pageDoc, err := goquery.NewDocumentFromReader(pageRes.Body)
			if err != nil {
				log.Fatal(err)
			}
			// Duyệt qua từng chương
			pageDoc.Find(".list-chapter li a").Each(func(index int, element *goquery.Selection) {
				// Lấy link chương
				link, exists := element.Attr("href")
				if !exists {
					return
				}

				// Lấy tên chương từ `title`
				title, _ := element.Attr("title")

				// Lưu vào danh sách
				chapters = append(chapters, map[string]string{
					"link":  link,
					"title": title,
				})
			})
		}
	} else {
		// Duyệt qua từng chương
		doc.Find(".list-chapter li a").Each(func(index int, element *goquery.Selection) {
			// Lấy link chương
			link, exists := element.Attr("href")
			if !exists {
				return
			}

			// Lấy tên chương từ `title`
			title, _ := element.Attr("title")

			// Lưu vào danh sách
			chapters = append(chapters, map[string]string{
				"link":  link,
				"title": title,
			})
		})
	}

	// Hiển thị kết quả
	for _, chapter := range chapters {
		fmt.Printf("Chương: %s\nLink: %s\n\n", chapter["title"], chapter["link"])
	}
	// Convert struct to JSON
	jsonData, err := json.MarshalIndent(storyInfo, "", "    ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
