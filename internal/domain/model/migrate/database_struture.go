package migrate

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/news"
	"awesomeProject/internal/domain/model/role"
	"awesomeProject/internal/domain/model/seeder"
	"awesomeProject/internal/domain/model/upload"
	"awesomeProject/internal/domain/model/user"
	"awesomeProject/internal/domain/model/user_deleted"
	"log"
)

func Migrate() {
	var count int64
	////
	TruncateTables()
	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}
	MigrateUser(count)
	////
	if err := database.DB.AutoMigrate(&role.Role{}); err != nil {
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

// /
func MigrateUser(count int64) {
	if err := database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}
	log.Println("Database models User migrated successfully")
	seeder.SeedUsers(10)
}

// /
func MigrateRole(count int64) {
	//eng: Drop roles table //ru: Полностью удаляем таблицу ролей
	if err := database.DB.Exec("DROP TABLE IF EXISTS roles_struct CASCADE;").Error; err != nil {
		log.Printf("Failed to drop roles table: %v", err)
	}

	//eng: Create roles table again with explicit types //ru: Создаем таблицу ролей заново с явным указанием типов
	if err := database.DB.Exec(`
		CREATE TABLE roles_struct (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			role_name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT NOT NULL
		);
	`).Error; err != nil {
		log.Fatalf("Failed to create roles table: %v", err)
	}

	database.DB.Model(&role.Role{}).Count(&count)
	log.Println("Database models Role migrated successfully")
	seeder.SeedRoles(3)
}

// /
func MigrateNews(count int64) {
	if err := database.DB.AutoMigrate(&news.News{}); err != nil {
		log.Fatalf("Failed to migrate news model: %v", err)
	}
	database.DB.Model(&news.News{}).Count(&count)
	log.Println("Database models News migrated successfully")
	seeder.SeedNews(5)
}

// /
func MigrateUserDeleted(count int64) {
	if err := database.DB.AutoMigrate(&user_deleted.UserDeleted{}); err != nil {
		log.Fatalf("Failed to migrate user_deleted model: %v", err)
	}
	database.DB.Model(&user_deleted.UserDeleted{}).Count(&count)
	log.Println("Database models UserDeleted migrated successfully")
	seeder.SeedUserDeleted(3)
}

func MigrateUpload(count int64) {
	if err := database.DB.AutoMigrate(&upload.Upload{}); err != nil {
		log.Fatalf("Failed to migrate upload model: %v", err)
	}
	database.DB.Model(&upload.Upload{}).Count(&count)
	log.Println("Database models Upload migrated successfully")
	seeder.SeedUpload(5)
}

// eng: TruncateTables clears all tables in the database //ru: TruncateTables очищает все таблицы в базе данных
func TruncateTables() {
	log.Println("Truncating all tables...")

	//eng: For PostgreSQL we use TRUNCATE //ru: Для PostgreSQL используем TRUNCATE
	if err := database.DB.Exec("TRUNCATE TABLE users_struct, roles_struct, news_struct, users_deleted_struct, uploads_struct CASCADE;").Error; err != nil {
		log.Printf("Failed to truncate tables: %v", err)
	}
	log.Println("All tables truncated successfully")
}
