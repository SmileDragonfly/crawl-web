package main

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestGetStoryByLink(t *testing.T) {
	type args struct {
		db  *gorm.DB
		url string
	}
	// Open db
	db, err := OpenDB("127.0.0.1", 3306, "truyenfull", "root", "123@123A")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetStoryByLink",
			args: args{
				db:  db,
				url: "https://truyenfull.vision/huou-con-va-vao-tim-tu-thoai-giang-son-bat-hieu/",
			},
			want: "https://truyenfull.vision/huou-con-va-vao-tim-tu-thoai-giang-son-bat-hieu/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStoryByLink(tt.args.db, tt.args.url); !reflect.DeepEqual(got.Url, tt.want) {
				t.Errorf("GetStoryByLink() = %v, want %v", got.Url, tt.want)
			}
		})
	}
}
