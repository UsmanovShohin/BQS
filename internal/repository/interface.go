package repository

import "queue-system/internal/models"

type IRepository interface {
	FindOrCreateByPhone(phone string) (*models.Client, error)
	Create(queue *models.QueueTicket) error
	GetLastByServiceType(serviceTypeId int) (*models.QueueTicket, error)
	FindNextWaitingByServiceType(serviceTypeId int) (*models.QueueTicket, error)
	FindActiveTicketByClient(clientID int) (*models.QueueTicket, error)
	Update(queue *models.QueueTicket) error
	GetCurrentClientByWindow(windowId int) (*models.QueueTicket, error)
	FindFreeByServiceType(serviceTypeID int) (*models.Window, error)
	GetWindowById(id int) (*models.Window, error)
	GetFreeWindowsSortedByLoad(serviceTypeID int) ([]models.WindowWithLoad, error)
	ListActiveCalls() ([]models.CallInfo, error)
	GetServiceTypeById(id int) (*models.ServiceType, error)
}
