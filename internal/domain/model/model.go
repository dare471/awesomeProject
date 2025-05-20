package model

import (
	"awesomeProject/internal/domain/model/migrate"
	"log"
)

// InitModels инициализирует все модели в базе данных
func InitModels() {
	log.Println("Initializing database models...")
	migrate.Migrate()
	log.Println("Database models initialized successfully")
}