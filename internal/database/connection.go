package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"queue-system/internal/config"
)

type Database struct {
	connection *gorm.DB
}

func NewDatabase(conn *gorm.DB) *Database {
	return &Database{connection: conn}
}

func ConnectDB(cfg *config.Configs) *gorm.DB {
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Username, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.DBName, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.SSLMode, cfg.DatabaseConfig.TimeZone,
	)

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных", err)
	}
	log.Println("Успешное подключение к базе данных")
	return db

}
