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
	Id           int
	UUID         string
	Username     string
	Email        string
	Password     string
	CreatedAt    time.Time
	StrCreatedAt string
	Role         string
	// IsMod        bool
	// IsAdmin      bool
}

type Comment struct {
	ID            int
	Post_ID       string
	User_UUID     string
	Username      string
	Text          string
	CreatedAt     time.Time
	Creation_Date string
	Creation_Hour string
	Like          int
	Dislike       int
}

type TemplateSetting struct {
	IsLogged bool
	Username string
}

type TemplateEdit struct {
	IsLogged   bool
	User       User
	Post       Post
	Categories []Category
}

type TemplateCreatedPost struct {
	IsLogged   bool
	User       User
	Categories []Category
}

type TemplateProfile struct {
	IsLogged bool
	User     User
	UserInfo UserInfo
}

type TemplatePost struct {
	IsLogged   bool
	User       User
	Category   string
	Posts      []Post
	Categories []Category
}

type TemplateComment struct {
	IsLogged bool
	User     User
	Post     Post
	Comments []Comment
}

type TemplateAdmin struct {
	IsLogged   bool
	User       User
	ReportInfo ReportInfo
	RequestMod []Request
}

type Request struct {
	User_id  int
	Username string
	Reason   string
}

type UserInfo struct {
	User          User
	PostedPost    []Post
	PostedComment []Comment
	LikedPost     []Post
	NbrLike       int
	NbrDislike    int
}

type Post2 struct {
	Id           int
	User_id      int
	Title        string
	Text         string
	NbrComments  int
	CreatedAt    time.Time
	StrCreatedAt string
	UpdatedAt    time.Time
	StrUpdatedAt string
	Likecount    int
	Dislikecount int
	IsLiked      bool
	IsDisliked   bool
	Username     string
	Categories   []Category
}

type Category struct {
	Id   int
	Name string
}

type ReportInfo struct {
	ReportPosts    []Post
	ReportComments []Comment
}
