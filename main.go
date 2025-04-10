package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	route "github.com/SinghaAnirban005/Lexi-Backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	db.AutoMigrate(&models.User{}, &models.Conversation{}, &models.Prompt{}, &models.Response{}, &models.Tag{}, &models.PromptTag{})

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173/",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	route.SetUpRoutes(app, db)

	log.Fatal(app.Listen(":8080"))
}
