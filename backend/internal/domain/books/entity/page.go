package entity

import "time"


type Page struct {
	ID         string
	BookID     string
	PageNumber int
	ImageURL   string
	CreatedAt  time.Time
}
