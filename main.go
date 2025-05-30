package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/delivery/http/router"
	"awesomeProject/internal/domain/model"
	_ "fmt"
	"log"
)

func main() {
	database.InitDatabase()
	model.InitModels()
	r := router.SetupRouter()
	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
