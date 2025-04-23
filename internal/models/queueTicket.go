package models

import "time"

type QueueTicket struct {
	ID            int        `json:"id"`
	Number        string     `json:"number"`
	ServiceTypeID int        `json:"service_type_id"`
	WindowID      *int       `json:"window_id"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	ServedAt      *time.Time `json:"served_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	ClientID      int        `json:"client_id"`
}

func (QueueTicket) TableName() string {
	return "queue_ticket"
}
