package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"queue-system/internal/models"
	"time"
)

func (s *Service) FinishClient(windowId int) (*models.QueueTicket, error) {
	queue, err := s.database.GetCurrentClientByWindow(windowId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("клиент не найден %v", err)
		}
		return nil, fmt.Errorf("ошибка при получении клиента %v", err)
	}
	if queue == nil {
		return nil, fmt.Errorf("у окна %d нет вызванного клиента", windowId)
	}
	now := time.Now()
	queue.Status = models.StatusFinished
	queue.FinishedAt = &now

	if err := s.database.Update(queue); err != nil {
		return nil, fmt.Errorf("ошибка при завершении очереди%v", err)
	}

	_, err = s.CallNextClient(windowId)
	if err != nil {
		log.Printf("Не удалось вызвать следующего клиента: %v", err)
	}
	return queue, nil
}
