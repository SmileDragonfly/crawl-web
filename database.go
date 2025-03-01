package main

import "time"

type tblType struct {
	ID          int
	Type        string
	Url         string
	CreatedTime time.Time
}

type tblStory struct {
	ID          int
	Title       string
	Image       string
	Author      string
	Type        string
	Status      string
	Source      string
	Rating      string
	RatingCount string
	Description string
	Url         string
	CreatedTime string
}

type tblChaper struct {
	ID      int
	StoryID int
	Title   string
	Url     string
	Content string
}
