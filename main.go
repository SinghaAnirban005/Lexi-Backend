package main

import (
	"log"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	route "github.com/SinghaAnirban005/Lexi-Backend/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgresql://neondb_owner:npg_8sUYPFLmvq1k@ep-divine-heart-a5mdcy8k-pooler.us-east-2.aws.neon.tech/neondb?sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	db.AutoMigrate(&models.User{}, &models.Conversation{}, &models.Prompt{}, &models.Response{}, &models.Tag{}, &models.PromptTag{})

	app := fiber.New()
	route.SetUpRoutes(app, db)

	log.Fatal(app.Listen(":8080"))
}
