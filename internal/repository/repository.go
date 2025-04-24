package repository

import (
	"errors"
	"gorm.io/gorm"
	"queue-system/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindOrCreateByPhone(phone string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("phone = ?", phone).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		client.Phone = phone
		err = r.db.Create(&client).Error
	}
	return &client, err
}

func (r *Repository) Create(queue *models.QueueTicket) error {
	return r.db.Create(queue).Error
}

func (r *Repository) GetLastByServiceType(serviceTypeId int) (*models.QueueTicket, error) {
	var queue models.QueueTicket
	err := r.db.Where("service_type_id = ?", serviceTypeId).
		Order("created_at desc").First(&queue).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &queue, err
}

func (r *Repository) FindNextWaitingByServiceType(serviceTypeId int) (*models.QueueTicket, error) {
	var queue models.QueueTicket
	err := r.db.Where("service_type_id = ? AND status = ?", serviceTypeId, "waiting").
		Order("created_at asc").First(&queue).Error
	return &queue, err
}

func (r *Repository) FindActiveTicketByClient(clientID int) (*models.QueueTicket, error) {
	var ticket models.QueueTicket
	err := r.db.Where("client_id = ? AND status IN ?", clientID, []string{"waiting", "called"}).
		Order("created_at desc").First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticket, err
}

func (r *Repository) Update(queue *models.QueueTicket) error {
	return r.db.Model(&models.QueueTicket{}).Where("id = ?", queue.ID).Updates(queue).Error
}

func (r *Repository) GetCurrentClientByWindow(windowId int) (*models.QueueTicket, error) {
	var queue models.QueueTicket
	err := r.db.Where("window_id = ? AND status = ? AND finished_at IS NULL", windowId, "called").
		Order("created_at asc").First(&queue).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &queue, err
}

func (r *Repository) FindFreeByServiceType(serviceTypeID int) (*models.Window, error) {
	var window models.Window
	err := r.db.Raw(`
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

func (r *Repository) GetWindowById(id int) (*models.Window, error) {
	var window models.Window
	err := r.db.First(&window, id).Error
	return &window, err
}

func (r *Repository) GetFreeWindowsSortedByLoad(serviceTypeID int) ([]models.WindowWithLoad, error) {
	var result []models.WindowWithLoad
	err := r.db.Table("windows").
		Select("windows.id, windows.name, COUNT(queue_ticket.id) as served_count").
		Joins("LEFT JOIN queue_ticket ON queue_ticket.window_id = windows.id AND queue_ticket.status = ?", models.StatusFinished).
		Where("windows.service_type_id = ?", serviceTypeID).
		Group("windows.id").
		Order("served_count asc").
		Scan(&result).Error
	return result, err
}

func (r *Repository) ListActiveCalls() ([]models.CallInfo, error) {
	var calls []models.CallInfo
	err := r.db.Table("queue_ticket").
		Select("queue_ticket.number AS ticket_number, windows.name AS window_name, queue_ticket.served_at").
		Joins("JOIN windows ON windows.id = queue_ticket.window_id").
		Where("queue_ticket.status = ?", models.StatusCalled).
		Order("queue_ticket.served_at asc").
		Scan(&calls).Error
	return calls, err
}

func (r *Repository) GetServiceTypeById(id int) (*models.ServiceType, error) {
	var serviceType models.ServiceType
	err := r.db.First(&serviceType, id).Error
	return &serviceType, err
}
