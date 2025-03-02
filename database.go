package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StoryInfo struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Author      string    `json:"author"`
	AuthorUrl   string    `json:"authorUrl"`
	Type        string    `json:"type"`
	Source      string    `json:"source"`
	Status      string    `json:"status"`
	Rate        string    `json:"rate"`
	RatingCount string    `json:"rating_count"`
	Description string    `json:"description"`
	CreatedTime time.Time `json:"created_time"`
}

type Chaper struct {
	ID      int
	StoryID int
	Title   string
	Url     string
	Content string
}

func OpenDB(dbIp string, dbPort int, dbName string, dbUser string, dbPassword string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbIp, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Tạo bảng nếu chưa tồn tại
	err = db.AutoMigrate(&StoryInfo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveStory(db *gorm.DB, story StoryInfo) error {
	result := db.Create(&story)
	return result.Error
}
