package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/delivery/Api"
	"awesomeProject/internal/domain/model"
	"awesomeProject/internal/domain/service"
	_ "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	database.InitDatabase()

	// Инициализация моделей
	model.InitModels()

	// Создаем сервис
	userService := service.NewUserService()

	r := gin.Default()

	// Endpoint для создания пользователя
	r.POST("/create/user", func(c *gin.Context) {
		var req service.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.CreateUser(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    user,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		response, err := userService.Login(creds.Email, creds.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	//
	protected := r.Group("/protected/user",
		Api.TokenAuthMiddleware())
	{
		protected.GET("/name/:id", func(c *gin.Context) {
			idParam := c.Param("id")
			id, err := strconv.ParseUint(idParam, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
				return
			}

			user, err := userService.GetUserByID(uint(id))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to retrieve user",
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"id":      idParam,
				"message": "Authorized",
				"data":    user,
			})
		})
		protected.GET("/all", func(c *gin.Context) {
			users, err := userService.GetAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "All users",
				"data":    users,
			})
		})
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World",
			})
		})
	}

	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
