package model

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/user"
	"log"
)

// InitModels инициализирует все модели в базе данных
func InitModels() {
	log.Println("Initializing database models...")

	// Инициализация модели пользователя
	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}

	// Проверяем, есть ли уже пользователи в базе
	var count int64
	database.DB.Model(&user.User{}).Count(&count)
	
	if count == 0 {
		log.Println("Creating initial admin user...")
		
		// Хешируем пароль
		hashedPassword, err := user.HashPassword("admin123")
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}

		// Создаем начального пользователя, если база пуста
		adminUser := user.User{
			Email:    "admin@example.com",
			Password: hashedPassword,
			Name:     "Admin",
			Age:      30,
			City:     "Moscow",
		}
		if err := database.DB.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		}
	}

	log.Println("Database models initialized successfully")
}


func Migrate() {
	database.DB.AutoMigrate(&user.User{})
	log.Println("Database models migrated successfully")
	database.DB.AutoMigrate(&news.News{})
	log.Println("Database models migrated successfully")
	
}
