package service

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"queue-system/internal/models"
	"queue-system/internal/repository"
	"queue-system/pkg"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	rep repository.IRepository
}

func NewService(rep repository.IRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) CreateTicketForClient(phone string, serviceTypeId int) (*models.QueueTicket, error) {
	if err := pkg.ValidatePhone(phone); err != nil {
		log.WithError(err).Warn("Неверный номер телефона")
		return nil, err
	}

	client, err := s.rep.FindOrCreateByPhone(phone)
	if err != nil {
		log.WithError(err).Error("Ошибка при поиске или создании клиента")
		return nil, err
	}
	log.WithField("client_id", client.ID).Info("Клиент найден или создан")

	existing, err := s.rep.FindActiveTicketByClient(client.ID)
	if err != nil {
		log.WithError(err).Error("Ошибка при проверке активного талона")
		return nil, err
	}
	if existing != nil {
		log.WithField("client_id", client.ID).Warn("Клиент уже в очереди")
		return existing, models.ErrAlreadyInQueue
	}

	queue, err := s.EnqueueClient(client.ID, serviceTypeId)
	if err != nil {
		log.WithError(err).Error("Не удалось создать талон в очереди")
		return nil, err
	}

	_ = s.TryAssignWindow(queue)
	return queue, nil
}

func (s *Service) EnqueueClient(clientID, serviceTypeID int) (*models.QueueTicket, error) {
	number, err := s.GenerateQueueNumber(serviceTypeID)
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

	err = s.rep.Create(queue)
	if err != nil {
		log.WithError(err).Error("Ошибка при добавлении талона в очередь")
	} else {
		log.WithFields(log.Fields{
			"client_id": clientID,
			"queue_num": number,
		}).Info("Талон успешно добавлен")
	}
	return queue, err
}

func (s *Service) GenerateQueueNumber(serviceTypeId int) (string, error) {
	serviceType, err := s.rep.GetServiceTypeById(serviceTypeId)
	if err != nil {
		return "", err
	}

	last, err := s.rep.GetLastByServiceType(serviceTypeId)
	if err != nil {
		return "", err
	}

	var number int
	if last != nil {
		numStr := strings.TrimPrefix(last.Number, serviceType.Code)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.WithError(err).Error("Ошибка преобразования номера очереди")
			return "", fmt.Errorf("неправильный формат номера: %v", err)
		}
		number = num + 1
	} else {
		number = 1
	}

	return fmt.Sprintf("%s%03d", serviceType.Code, number), nil
}

func (s *Service) TryAssignWindow(queue *models.QueueTicket) error {
	windows, err := s.rep.GetFreeWindowsSortedByLoad(queue.ServiceTypeID)
	if err != nil {
		log.WithError(err).Error("Не удалось получить свободные окна")
		return err
	}

	for _, w := range windows {
		current, err := s.rep.GetCurrentClientByWindow(w.ID)
		if err == nil && current == nil {
			queue.WindowID = &w.ID
			queue.Status = models.StatusCalled
			now := time.Now()
			queue.ServedAt = &now
			log.WithFields(log.Fields{
				"ticket": queue.Number,
				"window": w.ID,
			}).Info("Клиент направлен к окну")
			return s.rep.Update(queue)
		}
	}
	return nil
}

func (s *Service) AssignClientToWindow(queue *models.QueueTicket, windowId int) error {
	now := time.Now()
	queue.WindowID = &windowId
	queue.Status = models.StatusCalled
	queue.ServedAt = &now
	log.WithFields(log.Fields{
		"ticket": queue.Number,
		"window": windowId,
	}).Info("Клиент направлен к окну вручную")
	return s.rep.Update(queue)
}

func (s *Service) CallNextClient(windowId int) (*models.QueueTicket, error) {
	current, err := s.rep.GetCurrentClientByWindow(windowId)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении текущего клиента окна")
		return nil, err
	}
	if current != nil {
		log.WithField("window", windowId).Warn("Уже есть вызванный клиент")
		return nil, fmt.Errorf("у окна уже есть вызванный клиент")
	}

	window, err := s.rep.GetWindowById(windowId)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении окна")
		return nil, err
	}

	next, err := s.rep.FindNextWaitingByServiceType(window.ServiceTypeID)
	if err != nil || next == nil {
		log.WithError(err).Warn("Нет клиентов в очереди")
		return nil, errors.New("очередь пуста")
	}

	err = s.AssignClientToWindow(next, windowId)
	if err != nil {
		log.WithError(err).Error("Ошибка при назначении клиента на окно")
	}
	return next, err
}

func (s *Service) FinishClient(windowId int) (*models.QueueTicket, error) {
	client, err := s.rep.GetCurrentClientByWindow(windowId)
	if err != nil {
		log.WithError(err).Error("Ошибка при завершении — клиент не найден")
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("нет вызванного клиента у окна %d", windowId)
	}

	now := time.Now()
	client.Status = models.StatusFinished
	client.FinishedAt = &now

	if err := s.rep.Update(client); err != nil {
		log.WithError(err).Error("Ошибка при завершении клиента")
		return nil, err
	}

	log.WithFields(log.Fields{
		"ticket": client.Number,
		"window": windowId,
	}).Info("Клиент обслужен")

	_, err = s.CallNextClient(windowId)
	return client, nil
}

func (s *Service) GetCurrentClient(windowId int) (*models.QueueTicket, error) {
	return s.rep.GetCurrentClientByWindow(windowId)
}

func (s *Service) ListActiveCalls() ([]models.CallInfo, error) {
	return s.rep.ListActiveCalls()
}
