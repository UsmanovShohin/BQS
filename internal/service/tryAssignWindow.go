package service

import (
	"queue-system/internal/models"
	"time"
)

func (s *Service) TryAssignWindow(queue *models.QueueTicket) error {
	window, err := s.database.GetFreeWindowsSortedByLoad(queue.ServiceTypeID)
	if err != nil {
		return err
	}
	for _, window := range window {
		current, err := s.database.GetCurrentClientByWindow(window.ID)
		if err == nil && current == nil {
			queue.WindowID = &window.ID
			queue.Status = models.StatusCalled
			now := time.Now()
			queue.ServedAt = &now
			return s.database.Update(queue)
		}
	}
	return nil
}
