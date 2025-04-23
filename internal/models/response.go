package models

import (
	"errors"
	"time"
)

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

var ErrAlreadyInQueue = errors.New("клиент уже в очереди")
