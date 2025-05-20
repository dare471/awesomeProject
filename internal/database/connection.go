package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB // Подключение к БД через GORM

// InitDatabase создаёт соединение с базой данных
func InitDatabase() {
	log.Println("Initializing database connection...")
	
	dsn := "host=localhost user=laravel password=secret dbname=my_app_db port=5432 sslmode=disable TimeZone=UTC"
	log.Printf("Connecting to database with DSN: %s", dsn)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Проверяем соединение
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	DB = db
}

// GetDB возвращает экземпляр базы данных
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database connection is nil. Make sure InitDatabase() was called.")
	}
	return DB
}
