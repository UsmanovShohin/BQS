package service

import (
	"queue-system/internal/models"
	"time"
)

func (s *Service) EnqueueClient(clientID, serviceTypeID int) (*models.QueueTicket, error) {
	number, err := s.GeneraQueueNumber(serviceTypeID)
	if err != nil {
		return nil, err
	}

	queue := &models.QueueTicket{
		ClientID:      clientID,
		ServiceTypeID: serviceTypeID,
		Number:        number,
		Status:        models.StatusWaiting,
		CreatedAt:     time.Now(),
	}

	if err := s.database.Create(queue); err != nil {
		return nil, err
	}

	return queue, nil
}
