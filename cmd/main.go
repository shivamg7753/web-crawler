package main

import (
	"web-crawler/internal/crawler"
)

func main() {
	c := crawler.NewCrawler(10)
	c.Start("https://golang.org", 3)
}
