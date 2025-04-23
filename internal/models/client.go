package models

import "time"

type Client struct {
	ID        int       `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
