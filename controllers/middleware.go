package controller

import (
	"strings"

	"github.com/SinghaAnirban005/Lexi-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for token in multiple locations
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

// Helper function to extract token from multiple sources
func extractToken(c *fiber.Ctx) string {
	// 1. Check Authorization header (Bearer token)
	bearerToken := c.Get("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:6]) == "BEARER" {
		return bearerToken[7:]
	}

	// 2. Check query parameter
	token := c.Query("token")
	if token != "" {
		return token
	}

	// 3. Check cookie
	cookieToken := c.Cookies("token")
	if cookieToken != "" {
		return cookieToken
	}

	return ""
}
