package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"queue-system/internal/config"
)

type Database struct {
	connection *gorm.DB
}

func NewDatabase(conn *gorm.DB) *Database {
	return &Database{connection: conn}
}

func ConnectDB(cfg *config.Config) *gorm.DB {
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DB.Host, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port, cfg.DB.SSLMode, cfg.DB.TimeZone,
	)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных", err)
	}
	log.Println("Успешное подключение к базе данных")
	return db

}
