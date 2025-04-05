package controller

import (
	"net/http"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	"github.com/SinghaAnirban005/Lexi-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name,omitempty"`
	Username string `json:"username,omitempty"`
}

func SignUp(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload AuthPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		hashed, _ := utils.HashPassword(payload.Password)
		user := models.User{FullName: payload.FullName, Email: payload.Email, Username: payload.Username, Password: hashed}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "User creation failed"})
		}

		// token := utils.GenerateToken(user.ID)
		return c.JSON(fiber.Map{"user ID": user.ID})
	}
}

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload AuthPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		var user models.User
		db.Where("email= ?", payload.Email).First(&user)
		if user.ID == [16]byte{} || !utils.CheckPasswordHash(payload.Password, user.Password) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		token := utils.GenerateToken(user.ID.String())
		return c.JSON(fiber.Map{"token": token})
	}
}
