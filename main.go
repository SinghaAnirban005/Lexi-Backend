package main

import (
	"log"
	"os"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	route "github.com/SinghaAnirban005/Lexi-Backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	db.AutoMigrate(&models.User{}, &models.Conversation{}, &models.Prompt{}, &models.Response{}, &models.Tag{}, &models.PromptTag{})

	app := fiber.New()
	route.SetUpRoutes(app, db)

	log.Fatal(app.Listen(":8080"))
}
