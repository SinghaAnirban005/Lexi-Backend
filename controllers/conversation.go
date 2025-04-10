package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const LLMapiEndpoint = "https://api.groq.com/openai/v1/chat/completions"

func CreateConversation(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)

		var payload struct {
			Title string `json:"title"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		conv := models.Conversation{
			Title:   payload.Title,
			OwnerID: uuid.MustParse(userID),
		}
		db.Create(&conv)

		return c.JSON(conv)
	}
}

func GetUserConversations(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		var conversations []models.Conversation
		db.Where("owner_id = ?", userID).Preload("Prompts.Responses").Find(&conversations)
		return c.JSON(conversations)
	}
}

func CreatePromptWithResponse(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var payload struct {
			ConversationID string `json:"conversation_id"`
			PromptTitle    string `json:"prompt_title"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		prompt := models.Prompt{
			PromptTitle:    payload.PromptTitle,
			ConversationID: uuid.MustParse(payload.ConversationID),
		}
		if err := db.Create(&prompt).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create prompt",
			})
		}

		llmRequest := map[string]interface{}{
			"model": "llama-3.1-8b-instant",
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": []map[string]string{
						{"type": "text", "text": payload.PromptTitle},
					},
				},
			},
			"stream": false,
		}

		reqBody, err := json.Marshal(llmRequest)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to prepare AI request",
			})
		}

		req, err := http.NewRequest("POST", LLMapiEndpoint, bytes.NewBuffer(reqBody))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create AI request",
			})
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+"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")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to contact AI API",
			})
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read AI response",
			})
		}

		var llmResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error,omitempty"`
		}

		if err := json.Unmarshal(bodyBytes, &llmResp); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error":        "Failed to parse AI response",
				"raw_response": string(bodyBytes),
			})
		}

		if llmResp.Error != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "AI API error: " + llmResp.Error.Message,
			})
		}

		if len(llmResp.Choices) == 0 {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error":        "AI API returned no choices",
				"raw_response": string(bodyBytes),
			})
		}

		response := models.Response{
			Response: llmResp.Choices[0].Message.Content,
			PromptID: prompt.ID,
		}
		if err := db.Create(&response).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save response",
			})
		}

		return c.JSON(fiber.Map{
			"prompt":   prompt,
			"response": response,
		})
	}
}

func GetPromptsByConversation(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		conversationID := c.Params("conversation_id")
		if conversationID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Conversation ID is required",
			})
		}

		convUUID, err := uuid.Parse(conversationID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid conversation ID format",
			})
		}

		var conversation models.Conversation
		if err := db.Where("id = ? AND owner_id = ?", convUUID, userID).First(&conversation).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Conversation not found or access denied",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify conversation ownership",
			})
		}

		var prompts []models.Prompt
		if err := db.Where("conversation_id = ?", convUUID).
			Preload("Responses").
			Find(&prompts).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch prompts",
			})
		}

		return c.JSON(fiber.Map{
			"conversation_id": conversationID,
			"prompts":         prompts,
		})
	}
}
