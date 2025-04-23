package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"queue-system/internal/models"
)

func (r *Database) Create(queue *models.QueueTicket) error {
	if err := r.connection.Create(queue).Error; err != nil {
		return fmt.Errorf("create Queue Error: %v", err)
	}
	log.Printf("Успешно создание очереди %+v\n", queue)
	return nil

}

func (r *Database) GetLastByServiceType(serviceTypeId int) (*models.QueueTicket, error) {
	var queue models.QueueTicket
	err := r.connection.Where("service_type_id = ? ", serviceTypeId).
		Order("created_at desc").First(&queue).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last queue: %v", err)
	}
	return &queue, nil
}

func (r *Database) FindNextWaitingByServiceType(serviceTypeId int) (*models.QueueTicket, error) {
	var queue models.QueueTicket
	err := r.connection.Where("service_type_id = ? AND status = ?", serviceTypeId, "waiting").
		Order("created_at ASC").First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *Database) FindActiveTicketByClient(clientID int) (*models.QueueTicket, error) {
	var ticket models.QueueTicket
	err := r.connection.
		Where("client_id = ? AND status IN ?", clientID, []string{"waiting", "called"}).
		Order("created_at DESC").
		First(&ticket).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticket, err
}

func (r *Database) Update(queue *models.QueueTicket) error {
	return r.connection.Model(&models.QueueTicket{}).Where("id = ?", queue.ID).Updates(queue).Error
}
