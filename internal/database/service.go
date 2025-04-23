package database

import (
	"fmt"
	"queue-system/internal/models"
)

func (r *Database) GetServiceTypeById(id int) (*models.ServiceType, error) {
	var serviceType models.ServiceType
	if err := r.connection.First(&serviceType, id).Error; err != nil {
		return nil, fmt.Errorf("не найден тип услуги: %v", err)
	}
	return &serviceType, nil
}
