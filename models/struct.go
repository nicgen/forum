package models

import (
	"time"
)

type Post struct {
	ID            string
	Category_Name string
	User_UUID     string
	Category_ID   string
	Title         string
	Text          string
	ImagePath     string
	Like          int
	Dislike       int
	CreatedAt     time.Time
	Username      string
	Comments      []*Comment
	Status        string
	IsAuthor      string
	Data          map[string]interface{}
	Role          string
	ImageSize     int64
	Creation_Date string
	Creation_Hour string
}

type Notification struct {
	ID          int
	ReactionID  *int
	PostID      *int
	PostTitle string
	CommentID   *int
	CommentContent string
	IsOnComment bool // Indique si la réaction est sur un commentaire
	CreatedAt   time.Time
	Creation_Date string
	Creation_Hour string
	IsRead      bool
	TotalCount  *int // Total des notifications
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
	Data          map[string]interface{}
	Like          int
	Dislike       int
	Status        string
	IsAuthor      string
}

type Category struct {
	ID   int
	Name string
}

type Reports struct {
	ID              string
	User_UUID       string
	Username        string
	Post_ID         string
	Title           string
	Reported_Reason string
	Response_Text   string
}
