package article

import (
	"time"
)

type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Provider    string    `json:"provider"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pub_date"`
	Category    *string   `json:"category"`
	Thumbnail   *Image    `json:"thumbnail"`
}

type Image struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
