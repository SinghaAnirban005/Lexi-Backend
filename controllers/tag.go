package controller

import (
	"net/http"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTag(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tag models.Tag

		if err := c.BodyParser(&tag); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		if err := db.Create(&tag).Error; err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Tag creation failed"})
		}

		return c.JSON(tag)
	}
}

func AssignTagsToPrompt(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		promptID := c.Params("promptId")
		var payload struct {
			TagIDs []string `json:"tag_ids"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		for _, tagID := range payload.TagIDs {
			promptTag := models.PromptTag{PromptID: uuid.MustParse(promptID), TagID: uuid.MustParse(tagID)}
			db.FirstOrCreate(&promptTag)
		}
		return c.JSON(fiber.Map{"message": "Tags assigned successfully"})
	}
}
