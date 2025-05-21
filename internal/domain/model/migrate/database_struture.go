package migrate

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/user"
	"awesomeProject/internal/domain/model/news"
	"awesomeProject/internal/domain/model/user_deleted"
	"awesomeProject/internal/domain/model/upload"
	"log"
	"time"
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
	var count int64
	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}
	MigrateUser(count)
	////
	if err := database.DB.Automigrate(&role.Role{}); err != nil {
		log.Fatalf("Failed to migrate role model: %v", err)
	}
	MigrateRole(count)
	////
	if err := database.DB.AutoMigrate(&news.News{}); err != nil {
		log.Fatalf("Failed to migrate news model: %v", err)
	}
	MigrateNews(count)
	////
	if err := database.DB.AutoMigrate(&user_deleted.UserDeleted{}); err != nil {
		log.Fatalf("Failed to migrate user_deleted model: %v", err)
	}
	MigrateUserDeleted(count)
	////
	if err := database.DB.AutoMigrate(&upload.Upload{}); err != nil {
		log.Fatalf("Failed to migrate upload model: %v", err)
	}
	MigrateUpload(count)
}
///
func MigrateUser(count int64) {
	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}
	
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
///
func MigrateRole(count int64) {
	if err := database.DB.AutoMigrate(&role.Role{}); err != nil {
		log.Fatalf("Failed to migrate role model: %v", err)
	}
	database.DB.Model(&role.Role{}).Count(&count)
	if count == 0 {
		log.Println("Creating initial role...")
		role := role.Role{
			RoleName: "admin",
			Description: "Admin role",
		}
		if err := database.DB.Create(&role).Error; err != nil {
			log.Printf("Failed to create role: %v", err)
		}
	}
	log.Println("Database models Role migrated successfully")
}
///
func MigrateNews(count int64) {
	if err := database.DB.AutoMigrate(&news.News{}); err != nil {
		log.Fatalf("Failed to migrate news model: %v", err)
	}
	database.DB.Model(&news.News{}).Count(&count)
	if count == 0 {
		log.Println("Creating initial news...")
		news := news.News{
			Title: "Test News",
			Content: "This is a test news",
			Author: "Admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := database.DB.Create(&news).Error; err != nil {
			log.Printf("Failed to create news: %v", err)
		}
	}
	log.Println("Database models News migrated successfully")
}
///
func MigrateUserDeleted(count int64) {
	if err := database.DB.AutoMigrate(&user_deleted.UserDeleted{}); err != nil {
		log.Fatalf("Failed to migrate user_deleted model: %v", err)
	}
	database.DB.Model(&user_deleted.UserDeleted{}).Count(&count)
	if count == 0 {
		log.Println("Creating initial user_deleted...")
		user_deleted := user_deleted.UserDeleted{
			ID: 1,
			UserID: 1,
			DeletedAt: time.Now(),
		}
		if err := database.DB.Create(&user_deleted).Error; err != nil {
			log.Printf("Failed to create user_deleted: %v", err)
		}
	}
	log.Println("Database models UserDeleted migrated successfully")
}
///
func MigrateUpload(count int64) {
	if err := database.DB.AutoMigrate(&upload.Upload{}); err != nil {
		log.Fatalf("Failed to migrate upload model: %v", err)
	}
	database.DB.Model(&upload.Upload{}).Count(&count)
	if count == 0 {
		log.Println("Creating initial upload...")
		upload := upload.Upload{
			ID: 1,
			UserID: 1,
			Title: "Test Upload",
			Author: "Admin",
			Description: "This is a test upload",
			Content: "This is a test upload",
			Type: "text",
			Path: "test.txt",
		}
		if err := database.DB.Create(&upload).Error; err != nil {
			log.Printf("Failed to create upload: %v", err)
		}
	}
	log.Println("Database models Upload migrated successfully")
}