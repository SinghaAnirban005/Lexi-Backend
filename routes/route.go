package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUpRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/signup", Signup(db))
	app.Post("/login", Login(db))

	auth := app.Group("/api", AuthMiddleware())
	auth.Post("/conversations", CreateConversation(db))
	auth.Get("/conversations", GetUserConversations(db))
	auth.Post("/prompts", CreatePromptWithResponse(db))
	auth.Post("/tags", CreateTag(db))
	auth.Post("/prompts/:promptId/tags", AssignTagsToPrompt(db))
}
