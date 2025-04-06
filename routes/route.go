package route

import (
	controller "github.com/SinghaAnirban005/Lexi-Backend/controllers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUpRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/signup", controller.SignUp(db))
	app.Post("/login", controller.Login(db))

	auth := app.Group("/api", controller.AuthMiddleware())
	auth.Post("/conversations", controller.CreateConversation(db))
	auth.Get("/conversations", controller.GetUserConversations(db))
	auth.Post("/prompts", controller.CreatePromptWithResponse(db))
	auth.Post("/tags", controller.CreateTag(db))
	auth.Post("/prompts/:promptId/tags", controller.AssignTagsToPrompt(db))
}
