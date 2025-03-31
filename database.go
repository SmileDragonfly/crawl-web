package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// type StoryInfo struct {
// 	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"` // ID tự động tăng
// 	Title       string         `json:"title" gorm:"type:varchar(255);not null"`
// 	Image       string         `json:"image" gorm:"type:text"`
// 	Author      string         `json:"author" gorm:"type:varchar(255)"`
// 	AuthorUrl   string         `json:"authorUrl" gorm:"type:text"`
// 	Type        pq.StringArray `json:"type" gorm:"type:text[]"`
// 	Source      string         `json:"source" gorm:"type:varchar(255)"`
// 	Status      string         `json:"status" gorm:"type:varchar(50)"`
// 	Rate        string         `json:"rate" gorm:"type:varchar(10)"`
// 	RatingCount string         `json:"rating_count" gorm:"type:varchar(20)"`
// 	Description string         `json:"description" gorm:"type:text"`
// 	CreatedTime time.Time      `json:"created_time" gorm:"autoCreateTime"` // Tự động lấy timestamp khi tạo
// }

// type Chapter struct {
// 	ID            int       `gorm:"primaryKey;autoIncrement"`
// 	StoryID       int       `gorm:"not null;index"`
// 	StoryInfo     StoryInfo `gorm:"foreignKey:StoryID;references:ID;constraint:OnDelete:CASCADE"`
// 	ChapterNumber int       `gorm:"index"`
// 	Title         string    `gorm:"type:varchar(255);not null"`
// 	Url           string    `gorm:"type:text;"`
// 	Content       string    `gorm:"type:text"`
// 	CreatedTime   time.Time `gorm:"autoCreateTime"`
// }

type StoryInfo struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Url         string    `json:"url" gorm:"type:varchar(255);not null;unique"`
	Image       string    `json:"image" gorm:"type:text"`
	Author      string    `json:"author" gorm:"type:varchar(255)"`
	AuthorUrl   string    `json:"authorUrl" gorm:"type:text"`
	Type        string    `json:"type" gorm:"type:json"` // Chuyển từ pq.StringArray sang JSON
	Source      string    `json:"source" gorm:"type:varchar(255)"`
	Status      string    `json:"status" gorm:"type:varchar(50)"`
	Rate        string    `json:"rate" gorm:"type:varchar(10)"`
	RatingCount string    `json:"rating_count" gorm:"type:varchar(20)"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
}

type Chapter struct {
	ID            int       `gorm:"primaryKey;autoIncrement"`
	StoryID       int       `gorm:"not null;index"`
	StoryInfo     StoryInfo `gorm:"foreignKey:StoryID;references:ID;constraint:OnDelete:CASCADE"`
	ChapterNumber int       `gorm:"index"`
	Title         string    `gorm:"type:varchar(255);not null"`
	Url           string    `gorm:"type:varchar(255);not null;unique"`
	Content       string    `gorm:"type:text"`
	CreatedTime   time.Time `gorm:"autoCreateTime"`
}

func OpenDB(dbIp string, dbPort int, dbName string, dbUser string, dbPassword string) (*gorm.DB, error) {
	// dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	dbIp, dbPort, dbUser, dbPassword, dbName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbIp, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		return nil, err
	}
	// Tạo bảng nếu chưa tồn tại
	err = db.AutoMigrate(&StoryInfo{}, &Chapter{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveStory(db *gorm.DB, story *StoryInfo) error {
	result := db.Create(story)
	return result.Error
}

func SaveChapter(db *gorm.DB, chapter Chapter) error {
	result := db.Create(&chapter)
	return result.Error
}

func GetStoryByUrl(db *gorm.DB, url string) StoryInfo {
	var storyInfo StoryInfo
	db.Where("url = ?", url).First(&storyInfo)
	return storyInfo
}

func GetChapterByUrl(db *gorm.DB, url string) Chapter {
	var chapter Chapter
	db.Where("url = ?", url).First(&chapter)
	return chapter
}
