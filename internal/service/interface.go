package service

import "queue-system/internal/models"

type IService interface {
	EnqueueClient(clientID, serviceTypeID int) (*models.QueueTicket, error)
	ListActiveCalls() ([]models.CallInfo, error)
	CallNextClient(windowId int) (*models.QueueTicket, error)
	FinishClient(windowId int) (*models.QueueTicket, error)
	GetCurrentClient(windowId int) (*models.QueueTicket, error)
	TryAssignWindow(queue *models.QueueTicket) error
	AssignClientToWindow(queue *models.QueueTicket, windowId int) error
	CreateTicketForClient(phone string, serviceTypeId int) (*models.QueueTicket, error)
	GenerateQueueNumber(serviceTypeId int) (string, error)
}
