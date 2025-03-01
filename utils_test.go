package main

import "testing"

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
