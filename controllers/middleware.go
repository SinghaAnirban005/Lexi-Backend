package controller

import (
	"net/http"

	"github.com/SinghaAnirban005/Lexi-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		userID, err := utils.ParseToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		c.Locals("userID", userID)
		return c.Next()
	}
}
