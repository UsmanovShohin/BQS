package database

import "queue-system/internal/models"

type QueueRepository interface {
	Create(queue *models.QueueTicket) error
	GetLastByServiceType(serviceTypeId int) (*models.QueueTicket, error)
	FindNextWaitingByServiceType(serviceTypeId int) (*models.QueueTicket, error)
	Update(queue *models.QueueTicket) error
}
