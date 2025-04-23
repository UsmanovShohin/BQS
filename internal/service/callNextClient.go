package service

import (
	log "github.com/sirupsen/logrus"
	"queue-system/internal/database"
	"queue-system/internal/models"
	"time"
)

type Service struct {
	database *database.Database
}

func NewService(db *database.Database) *Service {
	return &Service{database: db}
}

func (s *Service) CallNextClient(windowId int) (*models.QueueTicket, error) {
	current, err := s.database.GetCurrentClientByWindow(windowId)
	if err != nil {
		log.WithError(err).Warn("Ошибка при получении текущего клиента из окна")
		return nil, err
	}
	if current != nil {
		log.WithError(err).Warn("У окна уже есть вызванный клиент")
		return nil, err
	}
	window, err := s.database.GetWindowById(windowId)
	if err != nil {
		log.WithError(err).Warn("Ошибка при получении ID окна")
		return nil, err
	}
	nextClient, err := s.database.FindNextWaitingByServiceType(window.ServiceTypeID)
	if err != nil || nextClient == nil {
		log.WithError(err).Warn("Нет подходящего клиента в очереди")
		return nil, err
	}

	nextClient.Status = models.StatusCalled
	nextClient.WindowID = &windowId
	now := time.Now()
	nextClient.ServedAt = &now

	err = s.database.Update(nextClient)
	if err != nil {
		log.WithError(err).Warn("Ошибка при обновлении данных")
		return nil, err
	}

	return nextClient, nil

}
