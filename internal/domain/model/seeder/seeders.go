package seeder

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/news"
	"awesomeProject/internal/domain/model/role"
	"awesomeProject/internal/domain/model/upload"
	"awesomeProject/internal/domain/model/user"
	"awesomeProject/internal/domain/model/user_deleted"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
)

func SeedUsers(count int) {
	log.Printf("Creating %d additional users...", count)

	var wg sync.WaitGroup
	users := make([]user.User, count)
	errors := make(chan error, count)

	// Создаем пользователей параллельно
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Создаем нового пользователя
			users[index] = user.User{
				Name:       faker.Name(),
				Age:        int(faker.RandomUnixTime() % 100),
				City:       faker.Word(),
				Email:      fmt.Sprintf("user%d@example.com", index+1),
				Role:       "user",
				IsActive:   true,
				IsVerified: true,
				IsDeleted:  false,
			}

			// Устанавливаем пароль
			hashedPassword, err := user.HashPassword("password123")
			if err != nil {
				errors <- fmt.Errorf("failed to hash password for user %d: %v", index, err)
				return
			}
			users[index].Password = hashedPassword

			// Создаем пользователя в базе данных
			if err := database.DB.Create(&users[index]).Error; err != nil {
				errors <- fmt.Errorf("failed to create user %s: %v", users[index].Email, err)
				return
			}
		}(i)
	}

	// Запускаем горутину для обработки ошибок
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Обрабатываем ошибки
	for err := range errors {
		log.Printf("Error creating user: %v", err)
	}

	log.Printf("Successfully created users")
}

func SeedRoles(count int) {
	log.Printf("Creating %d additional roles...", count)

	roles := make([]role.Role, count)

	for i := 0; i < count; i++ {
		roles[i] = role.Role{
			RoleName:    fmt.Sprintf("role%d", i+1),
			Description: faker.Sentence(),
		}

		if err := database.DB.Create(&roles[i]).Error; err != nil {
			log.Printf("Failed to create role %s: %v", roles[i].RoleName, err)
		}
	}

	log.Printf("Successfully created roles")
}

func SeedNews(count int) {
	log.Printf("Creating %d additional news...", count)

	newsItems := make([]news.News, count)

	for i := 0; i < count; i++ {
		// limints length of title to 255 characters
		// Ограничиваем все строковые поля до 255 символов

		title := faker.Sentence()
		if len(title) > 255 {
			title = title[:255]
		}

		description := faker.Sentence()
		if len(description) > 255 {
			description = description[:255]
		}

		content := faker.Sentence() //eng: Use a Sentence for Paragraph to limit length //ru: Используем Sentence вместо Paragraph для ограничения длины

		if len(content) > 255 {
			content = content[:255]
		}

		author := faker.Name()
		if len(author) > 255 {
			author = author[:255]
		}

		category := faker.Word() // Used Word for Category // ru: Используем Word для категории
		if len(category) > 255 {
			category = category[:255]
		}

		image := fmt.Sprintf("image_%d.jpg", i+1) // Generate a name of image file // ru: Генерируем имя файла изображения
		if len(image) > 255 {
			image = image[:255]
		}

		newsItems[i] = news.News{
			Title:       title,
			Description: description,
			Content:     content,
			Author:      author,
			Category:    category,
			Image:       image,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := database.DB.Create(&newsItems[i]).Error; err != nil {
			log.Printf("Failed to create news %s: %v", newsItems[i].Title, err)
		}
	}

	log.Printf("Successfully created news")
}

func SeedUserDeleted(count int) {
	log.Printf("Creating %d additional user_deleted...", count)

	userDeleted := make([]user_deleted.UserDeleted, count)

	for i := 0; i < count; i++ {
		// Create a unique user ID for each deleted user
		userID := uint(i + 1)

		userDeleted[i] = user_deleted.UserDeleted{
			UserID:    userID,
			DeletedAt: time.Now(),
			User: user.User{
				Name:          faker.Name(),
				Age:           int(faker.RandomUnixTime() % 100),
				City:          faker.Word(),
				Email:         fmt.Sprintf("deleted_user_%d@example.com", i+1),
				Password:      "deleted_user_password", // Inactive users don't need secure passwords
				Role:          "deleted_user",
				IsActive:      false,
				IsActive_at:   time.Now(),
				IsVerified:    false,
				IsVerified_at: time.Now(),
				IsDeleted:     true,
			},
		}

		if err := database.DB.Create(&userDeleted[i]).Error; err != nil {
			log.Printf("Failed to create user_deleted %d: %v", i+1, err)
			continue
		}
		log.Printf("Successfully created user_deleted %d", i+1)
	}
}

func SeedUpload(count int) {
	log.Printf("Creating %d additional uploads...", count)
	uploads := make([]upload.Upload, count)

	for i := 0; i < count; i++ {
		// Limints length of description to 200 characters // ru: Ограничиваем длину описания до 200 символов
		title := faker.Sentence()
		if len(title) > 200 {
			title = title[:200]
		}
		// Limints length of description to 200 characters // ru: Ограничиваем длину описания до 200 символов
		description := faker.Sentence()
		if len(description) > 200 {
			description = description[:200]
		}

		// Limints length of description to 1000 characters // ru: Ограничиваем длину контента до 1000 символов
		content := faker.Paragraph()
		if len(content) > 200 {
			content = content[:200]
		}

		uploads[i] = upload.Upload{
			Title:       title,
			Content:     content,
			Author:      faker.Name(),
			Description: description,
			Type:        "text",
			Path:        fmt.Sprintf("file_%d.txt", i+1),
		}

		if err := database.DB.Create(&uploads[i]).Error; err != nil {
			log.Printf("Failed to create upload %s: %v", uploads[i].Title, err)
		}
	}
	log.Printf("Successfully created uploads")
}
