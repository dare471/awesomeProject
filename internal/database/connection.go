package database

import (
	"fmt"
	"log"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Подключение к БД через GORM

// getEnvWithDefault returns the value of the environment variable or a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDatabase создаёт соединение с базой данных
func InitDatabase() {
	log.Println("Initializing database connection...")
	host := getEnvWithDefault("DB_HOST", "localhost")
	user := getEnvWithDefault("DB_USER", "laravel")
	password := getEnvWithDefault("DB_PASSWORD", "secret")
	dbname := getEnvWithDefault("DB_NAME", "my_app_db")
	port := getEnvWithDefault("DB_PORT", "5432")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)
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
