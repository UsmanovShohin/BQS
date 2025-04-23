package models

import "time"

type Window struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ServiceTypeID int    `json:"service_type_id"`
}

type WindowWithLoad struct {
	ID          int
	Name        string
	ServedCount int
}

type CallInfo struct {
	TicketNumber string    `json:"number"`
	WindowName   string    `json:"window"`
	CalledAt     time.Time `json:"called_at"`
}
