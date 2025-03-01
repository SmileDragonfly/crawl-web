package main

type Config struct {
	// Database
	DBIP       string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	// Source url
	SourceUrl string
}
