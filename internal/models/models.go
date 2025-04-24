package models

import "time"

type Admin struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type Client struct {
	ID        int       `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Code    int     `json:"code"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Client  *Client `json:"client"`
	Data    Data    `json:"data"`
}

type Data struct {
	Id            int        `json:"id"`
	Number        string     `json:"number"`
	Phone         string     `json:"phone"`
	ServiceTypeId int        `json:"service_type_id"`
	WindowId      *int       `json:"window_id,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	ServedAt      *time.Time `json:"served_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

type Request struct {
	WindowId      int `json:"window_id"`
	ServiceTypeId int `json:"service_type_id"`
}

type RegisterAndQueueRequest struct {
	Phone         string `json:"phone"`
	ServiceTypeId int    `json:"service_type_id"`
}

type ServiceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

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

// QueueTicket - DB model
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
