package models

import (
	"time"
)

type Post struct {
	ID            string
	Category_ID   string
	User_UUID     string
	Username      string
	Title         string
	Text          string
	Status        string
	Like          int
	Dislike       int
	Data          map[string]interface{}
	Comments      []*Comment
	CreatedAt     time.Time
	Creation_Date string
	Creation_Hour string
}

type User struct {
	Id           string
	UUID         string
	Username     string
	Email        string
	Password     string
	CreatedAt    time.Time
	StrCreatedAt string
	Role         string
	IsRequest    bool
}

type Comment struct {
	ID            string
	Post_ID       string
	User_UUID     string
	Username      string
	Text          string
	CreatedAt     time.Time
	Creation_Date string
	Creation_Hour string
	Like          int
	Dislike       int
	Status        string
}

type Category struct {
	ID   int
	Name string
}
