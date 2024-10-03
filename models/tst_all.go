package models

// PageData: data structure passed to templates, it holds the common data for all pages
type PageData struct {
	Title     string
	Header    string
	Content   interface{}
	IsError   bool
	ErrorCode int
}
