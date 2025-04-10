package controller

import (
	"log"
	"strings"

	"github.com/SinghaAnirban005/Lexi-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractToken(c)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
		}

		userID, err := utils.ParseToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) string {
	log.Println("All headers:", c.GetReqHeaders())
	// log.Println("All cookies:", c.Cookies())

	bearerToken := c.Get("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:6]) == "BEARER" {
		return bearerToken[7:]
	}

	token := c.Query("token")
	if token != "" {
		return token
	}

	cookieToken := c.Cookies("token")
	if cookieToken != "" {
		return cookieToken
	}

	log.Println("Cookiees --> ", cookieToken)

	return ""
}
