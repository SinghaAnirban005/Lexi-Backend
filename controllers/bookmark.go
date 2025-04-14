package controller

import (
	"github.com/SinghaAnirban005/Lexi-Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBookmark(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)

		var payload struct {
			ConversationID string `json:"conversation_id"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot Parse JSON",
			})
		}

		bookmark := models.Bookmark{
			UserID:         uuid.MustParse(userID),
			ConversationID: uuid.MustParse(payload.ConversationID),
		}

		db.Create(&bookmark)

		return c.JSON(bookmark)
	}
}
