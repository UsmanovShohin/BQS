package database

import (
	"errors"
	"gorm.io/gorm"
	"queue-system/internal/models"
)

func (r *Database) FindOrCreateByPhone(phone string) (*models.Client, error) {
	var client models.Client
	err := r.connection.Where("phone = ?", phone).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		client = models.Client{
			Phone: phone}
		err = r.connection.Create(&client).Error
	}
	return &client, err

}
