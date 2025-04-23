package service

import (
	"queue-system/internal/models"
	"time"
)

func (s *Service) AssignClientToWindow(queue *models.QueueTicket, windowId int) error {
	now := time.Now()
	queue.Status = models.StatusCalled
	queue.WindowID = &windowId
	queue.ServedAt = &now
	return s.database.Update(queue)
}
