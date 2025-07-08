package routes

import (
	"tether-server/handlers"
	"tether-server/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Tether Messenger API is running",
		})
	})

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Post("/request-code", handlers.RequestCode)
	auth.Post("/verify-code", handlers.VerifyCode)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware())
	protected.Get("/chats", handlers.GetChats)
	protected.Post("/chats", handlers.CreateChat)
	protected.Get("/chats/:chatId/messages", handlers.GetMessages)
	protected.Post("/messages", handlers.SendMessage)
	protected.Get("/users/search", handlers.SearchUsers)
	protected.Get("/profile", handlers.GetProfile)
	protected.Put("/profile", handlers.UpdateProfile)
	protected.Post("/profile/avatar", handlers.UploadAvatar)
}
