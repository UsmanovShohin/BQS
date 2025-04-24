package models

import "errors"

const (
	StatusWaiting  = "waiting"
	StatusCalled   = "called"
	StatusFinished = "finished"
)

var ErrAlreadyInQueue = errors.New("клиент уже в очереди")
