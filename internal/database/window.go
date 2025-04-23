package database

import (
	"errors"
	"gorm.io/gorm"
	"queue-system/internal/models"
)

func (r *Database) GetCurrentClientByWindow(windowId int) (*models.QueueTicket, error) {
	var queue models.QueueTicket
	err := r.connection.Where("window_id = ? AND status = ? AND finished_at IS NULL", windowId, "called").
		Order("created_at ASC").First(&queue).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *Database) FindFreeByServiceType(serviceTypeID int) (*models.Window, error) {
	var window models.Window
	err := r.connection.Raw(`
		SELECT * FROM windows w
		WHERE w.service_type_id = ?
		AND NOT EXISTS (
			SELECT 1 FROM queue_ticket q
			WHERE q.window_id = w.id AND q.status = 'called' AND q.finished_at IS NULL
		)
		LIMIT 1;
	`, serviceTypeID).Scan(&window).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &window, err
}

func (r *Database) GetWindowById(id int) (*models.Window, error) {
	var window models.Window
	err := r.connection.First(&window, id).Error
	if err != nil {
		return nil, err
	}
	return &window, nil
}

func (r *Database) GetFreeWindowsSortedByLoad(serviceTypeID int) ([]models.WindowWithLoad, error) {
	var result []models.WindowWithLoad

	// Здесь мы выбираем окна по типу услуги и считаем сколько клиентов они обслужили (finished)
	err := r.connection.
		Table("windows").
		Select("windows.id, windows.name, COUNT(queue_ticket.id) as served_count").
		Joins("LEFT JOIN queue_ticket ON queue_ticket.window_id = windows.id AND queue_ticket.status = ?", models.StatusFinished).
		Where("windows.service_type_id = ?", serviceTypeID).
		Group("windows.id").
		Order("served_count ASC").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Database) ListActiveCalls() ([]models.CallInfo, error) {
	var res []models.CallInfo
	err := r.connection.
		Table("queue_ticket").
		Select("queue_ticket.number AS ticket_number, windows.name AS window_name, queue_ticket.served_at").
		Joins("JOIN windows ON windows.id = queue_ticket.window_id").
		Where("queue_ticket.status = ?", models.StatusCalled).
		Order("queue_ticket.served_at ASC").
		Scan(&res).Error
	return res, err
}
