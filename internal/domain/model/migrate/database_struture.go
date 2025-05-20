package migrate

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/user"
	"awesomeProject/internal/domain/model/news"
	"awesomeProject/internal/domain/model/user_deleted"
	"awesomeProject/internal/domain/model/upload"
	"log"
)

func Migrate() {
	// Сначала добавляем поле role с значением по умолчанию
	if err := database.DB.Exec("ALTER TABLE users_struct ADD COLUMN IF NOT EXISTS role varchar(255) DEFAULT 'user'").Error; err != nil {
		log.Fatalf("Failed to add role column: %v", err)
	}
	
	// Затем делаем поле обязательным
	if err := database.DB.Exec("ALTER TABLE users_struct ALTER COLUMN role SET NOT NULL").Error; err != nil {
		log.Fatalf("Failed to set role column as NOT NULL: %v", err)
	}

	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}
	MigrateUser()
	if err := database.DB.AutoMigrate(&news.News{}); err != nil {
		log.Fatalf("Failed to migrate news model: %v", err)
	}
	MigrateNews()
	if err := database.DB.AutoMigrate(&user_deleted.UserDeleted{}); err != nil {
		log.Fatalf("Failed to migrate user_deleted model: %v", err)
	}
	MigrateUserDeleted()
	if err := database.DB.AutoMigrate(&upload.Upload{}); err != nil {
		log.Fatalf("Failed to migrate upload model: %v", err)
	}
	MigrateUpload()
}

func MigrateUser() {
	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}
	var count int64
	database.DB.Model(&user.User{}).Count(&count)
	
	if count == 0 {
		log.Println("Creating initial admin user...")
		
		hashedPassword, err := user.HashPassword("admin123")
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}
		
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
	log.Println("Database models User migrated successfully")
}

func MigrateNews() {
	if err := database.DB.AutoMigrate(&news.News{}); err != nil {
		log.Fatalf("Failed to migrate news model: %v", err)
	}
	var count int64
	database.DB.Model(&news.News{}).Count(&count)
	if count == 0 {
		log.Println("Creating initial news...")
		
	}
	log.Println("Database models News migrated successfully")
}

func MigrateUserDeleted() {
	if err := database.DB.AutoMigrate(&user_deleted.UserDeleted{}); err != nil {
		log.Fatalf("Failed to migrate user_deleted model: %v", err)
	}
	log.Println("Database models UserDeleted migrated successfully")
}

func MigrateUpload() {
	if err := database.DB.AutoMigrate(&upload.Upload{}); err != nil {
		log.Fatalf("Failed to migrate upload model: %v", err)
	}
	log.Println("Database models Upload migrated successfully")
}