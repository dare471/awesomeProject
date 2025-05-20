package Api

import (
	"awesomeProject/internal/domain/service"
	"awesomeProject/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthMiddleware struct {
	authService usecase.AuthService
}

func NewAuthMiddleWare(authService usecase.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	// Создаем сервис при каждом вызове middleware
	userService := service.NewUserService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("Received Authorization header: %s", authHeader)
		
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Проверяем формат Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Invalid Authorization header format: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		log.Printf("Processing token: %s", tokenString)
		
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Парсим и валидируем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(service.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Проверяем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Проверяем наличие user_id в claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// Проверяем срок действия токена
		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid expiration time in token"})
			c.Abort()
			return
		}
		if float64(time.Now().Unix()) > exp {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// После успешного получения userID
		log.Printf("Extracted userID from token: %d", userID)
		
		user, err := userService.GetUserByID(userID)
		if err != nil {
			log.Printf("Error getting user by ID: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User not found: %v", err)})
			c.Abort()
			return
		}
		
		log.Printf("Successfully authenticated user: %+v", user)
		c.Set("user", user)
		c.Next()
	}
}


func (a *AuthMiddleware) TokenValidatorMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		_, err := a.authService.Authenticate(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (a *AuthMiddleware) StatusCodeBadRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
		})
		c.Abort()
	}
}
