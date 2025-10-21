package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"web-crawler/internal/db"
	"web-crawler/internal/model"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	visited   map[string]bool
	mu        sync.Mutex
	store     *db.Database
	jobs      chan string
	results   chan model.Page
	wg        sync.WaitGroup
	jobsWg    sync.WaitGroup
	maxPages  int
	pageCount int
}

func NewCrawler(maxPages int) *Crawler {
	store := db.InitDB("./crawler.db")
	return &Crawler{
		visited:  make(map[string]bool),
		store:    store,
		jobs:     make(chan string, 100),
		results:  make(chan model.Page, 100),
		maxPages: maxPages,
	}
}

func (c *Crawler) Crawl(url string) (model.Page, []string) {
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf(" Error fetching %s: %v\n", url, err)
		return model.Page{URL: url}, nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return model.Page{URL: url}, nil
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	text := doc.Find("body").Text()
	wordCount := len(strings.Fields(text))

	var links []string
	baseURL := "https://golang.org"

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}
		var fullURL string
		if strings.HasPrefix(href, "http") {
			fullURL = href
		} else if strings.HasPrefix(href, "/") {
			fullURL = baseURL + href
		}
		if fullURL != "" && strings.Contains(fullURL, "golang.org") {
			links = append(links, fullURL)
		}
	})

	page := model.Page{
		URL:       url,
		Title:     title,
		WordCount: wordCount,
		LinkCount: len(links),
	}
	return page, links
}

func (c *Crawler) Worker() {
	defer c.wg.Done()
	for url := range c.jobs {
		page, links := c.Crawl(url)
		c.results <- page
		for _, link := range links {
			c.mu.Lock()
			if !c.visited[link] && c.pageCount < c.maxPages {
				c.visited[link] = true
				c.pageCount++
				c.jobsWg.Add(1)
				go func(l string) { c.jobs <- l }(link)
			}
			c.mu.Unlock()
		}
		c.jobsWg.Done()
	}
}

func (c *Crawler) Start(startURL string, workers int) {
	fmt.Println("Starting crawler...")
	fmt.Printf("Max pages: %d | Workers: %d\n\n", c.maxPages, workers)

	c.visited[startURL] = true
	c.pageCount = 1
	c.jobsWg.Add(1)
	c.jobs <- startURL

	for i := 0; i < workers; i++ {
		c.wg.Add(1)
		go c.Worker()
	}

	go func() {
		c.jobsWg.Wait()
		close(c.jobs)
	}()
	go func() {
		c.wg.Wait()
		close(c.results)
	}()

	count := 0
	for page := range c.results {
		count++
		c.store.SavePage(page)
		fmt.Printf("✓ [%d/%d] %s\n   Title: %s | Words: %d | Links: %d\n\n",
			count, c.maxPages, page.URL, page.Title, page.WordCount, page.LinkCount)
	}

	c.store.Close()
	fmt.Println("Done! Check crawler.db")
	fmt.Printf("Total pages crawled: %d\n", count)
}
