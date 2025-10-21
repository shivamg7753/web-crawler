# Go Web Crawler

A **high-performance concurrent web crawler** built in **Go** that scrapes web pages, extracts their titles, counts words and links, and stores everything in a local **SQLite** database.

Built with using **Goroutines**, **Channels**, and **WaitGroups** for clean concurrency control.

---

## Features

Multi-threaded crawling using Go routines  
Smart link discovery and deduplication  
Real-time progress logs in the terminal  
SQLite3 database integration for data storage  
Graceful shutdown and concurrency-safe design  
Modular architecture (`cmd`, `internal/db`, `internal/crawler`, `internal/model`)


## Tech Stack

| Component       | Description                                             |
| --------------- | ------------------------------------------------------- |
| **Language**    | Go 1.20+                                                |
| **Database**    | SQLite3                                                 |
| **HTML Parser** | [goquery](https://github.com/PuerkitoBio/goquery)       |
| **Concurrency** | Goroutines + Channels + sync.WaitGroup                  |
| **Driver**      | [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) |

---

## Installation & Setup

### 1️ Clone the repository

```bash
git clone https://github.com/shivamg7753/web-crawler.git
cd web-crawler

Install dependencies
go mod tidy

Run the crawler
go run ./cmd/main.go


Example Output


Starting crawler...
Max pages: 10 | Workers: 3

✓ [1/10] https://golang.org
   Title: The Go Programming Language | Words: 2250 | Links: 34

✓ [2/10] https://golang.org/doc/
   Title: Documentation - The Go Programming Language | Words: 1020 | Links: 12

Done! Check crawler.db
Total pages crawled: 10
```
