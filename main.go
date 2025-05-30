package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/delivery/http/router"
	"awesomeProject/internal/domain/model"
	_ "fmt"
)

func main() {
	database.InitDatabase()
	model.InitModels()
	router.SetupRouter()
}
