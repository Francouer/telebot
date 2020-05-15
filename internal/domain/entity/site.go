package entity

import "time"

//Site - contains all site entity information (ID, URL, Name, RequestTimeout, ResponseStatus, Description,
// CreatedAt)
type Site struct {
	ID             uint
	UserID         uint
	Name           string
	URL            string
	RequestTimeout int64
	ResponseStatus int64
	Description    string
	CreatedAt      time.Time
}
