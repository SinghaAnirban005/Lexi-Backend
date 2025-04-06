package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/SinghaAnirban005/Lexi-Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const LLMapiEndpoint = "https://api.groq.com/openai/v1/chat/completions"

func CreateConversation(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		title := c.FormValue("title")
		conv := models.Conversation{Title: title, OwnerID: uuid.MustParse(userID)}
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
		prompt := models.Prompt{PromptTitle: payload.PromptTitle, ConversationID: uuid.MustParse(payload.ConversationID)}
		db.Create(&prompt)

		llmRequest := map[string]interface{}{
			"model": "llama3-70b-8192",
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": []map[string]string{
						{
							"type": "text",
							"text": payload.PromptTitle,
						},
					},
				},
			},
			"stream": false,
		}

		reqBody, _ := json.Marshal(llmRequest)
		req, _ := http.NewRequest("POST", LLMapiEndpoint, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LLM_API_KEY"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to contact AI API"})
		}
		defer resp.Body.Close()
		fmt.Println(os.Getenv("LLM_API_KEY"))
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		var llmResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(bodyBytes, &llmResp); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse AI response"})
		}

		generated := llmResp.Choices[0].Message.Content

		response := models.Response{
			Response: generated,
			PromptID: prompt.ID,
		}
		db.Create(&response)

		return c.JSON(fiber.Map{"prompt": prompt, "response": response})
	}
}
