package service

import "queue-system/internal/models"

func (s *Service) GetCurrentClient(windowId int) (*models.QueueTicket, error) {
	return s.database.GetCurrentClientByWindow(windowId)
}
