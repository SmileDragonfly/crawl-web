package main

import (
	"os"
	"strconv"
	"testing"
)

func TestGetStoryByChapterNumber(t *testing.T) {
	type args struct {
		sourceUrl string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetStoryByChapterNumber",
			args: args{
				sourceUrl: "https://truyenfull.vision",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetStoryByChapterNumber(tt.args.sourceUrl)
		})
	}
}

func TestWriteToString(t *testing.T) {
	f, err := os.Create("text.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	for i := range 10 {
		f.WriteString(strconv.Itoa(i))
	}
}

func TestOpenDB(t *testing.T) {
	type args struct {
		dbIp       string
		dbPort     int
		dbName     string
		dbUser     string
		dbPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestOpenDB",
			args: args{
				dbIp:       "127.0.0.1",
				dbPort:     5432,
				dbName:     "truyenfull",
				dbUser:     "postgres",
				dbPassword: "123@123A",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OpenDB(tt.args.dbIp, tt.args.dbPort, tt.args.dbName, tt.args.dbUser, tt.args.dbPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetStoryInfoFromFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetStoryInfoFromFile",
			args: args{
				filePath: "StoryLink.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetStoryInfoFromFile(tt.args.filePath)
		})
	}
}
