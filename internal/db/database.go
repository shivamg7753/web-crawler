package db

import (
	"database/sql"
	"fmt"
	"web-crawler/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func InitDB(path string) *Database {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pages (
		url TEXT PRIMARY KEY,
		title TEXT,
		word_count INTEGER,
		link_count INTEGER,
		crawled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}

	return &Database{DB: db}
}

func (d *Database) SavePage(page model.Page) {
	_, err := d.DB.Exec(`INSERT OR REPLACE INTO pages (url, title, word_count, link_count) 
		VALUES (?, ?, ?, ?)`,
		page.URL, page.Title, page.WordCount, page.LinkCount)
	if err != nil {
		fmt.Printf("DB Error: %v\n", err)
	}
}

func (d *Database) Close() {
	d.DB.Close()
}
