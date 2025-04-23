package service

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"queue-system/internal/models"
	"regexp"
)

func (s *Service) CreateTicketForClient(phone string, serviceTypeId int) (*models.QueueTicket, error) {
	if err := s.ValidatePhone(phone); err != nil {
		log.WithError(err).Warn("Некорректный номер телефона")
		return nil, err
	}

	client, err := s.database.FindOrCreateByPhone(phone)
	if err != nil {
		log.WithError(err).Error("Ошибка при поиске/создании клиента")
		return nil, err
	}
	log.WithField("clientID", client.ID).Info("Клиент найден или создан")

	existTicket, err := s.database.FindActiveTicketByClient(client.ID)
	if err != nil {
		log.WithError(err).Error("Ошибка при проверке активного талона")
		return nil, err
	}

	if existTicket != nil {
		return existTicket, models.ErrAlreadyInQueue
	}

	queue, err := s.EnqueueClient(client.ID, serviceTypeId)
	if err != nil {
		log.WithError(err).Error("Ошибка при создании талона")
		return nil, err
	}
	_ = s.TryAssignWindow(queue)

	return queue, nil

}

func (s *Service) ValidatePhone(phone string) error {
	if len(phone) > 9 {
		return errors.New("номер телефона не должен превышать 9 цифр")
	}
	matched, err := regexp.MatchString(`^\d{1,9}$`, phone)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("номер телефона должен содержать только цифры")
	}
	return nil
}
