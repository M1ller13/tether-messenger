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
	auth.Post("/verify-email", handlers.VerifyEmail)
	auth.Post("/request-password-reset", handlers.RequestPasswordReset)
	auth.Post("/reset-password", handlers.ResetPassword)
	auth.Post("/refresh-token", handlers.RefreshToken)
	auth.Post("/logout", handlers.Logout)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware())

	// Chat routes
	protected.Get("/chats", handlers.GetChats)
	protected.Post("/chats", handlers.CreateChat)
	protected.Get("/chats/:chatId/messages", handlers.GetMessages)
	protected.Post("/messages", handlers.SendMessage)

	// User routes
	protected.Get("/users/search", handlers.SearchUsers)
	protected.Get("/profile", handlers.GetProfile)
	protected.Put("/profile", handlers.UpdateProfile)
	protected.Post("/profile/avatar", handlers.UploadAvatar)

	// Board routes
	protected.Get("/boards", handlers.GetBoards)
	protected.Post("/boards", handlers.CreateBoard)
	protected.Get("/boards/:id", handlers.GetBoard)
	protected.Put("/boards/:id", handlers.UpdateBoard)
	protected.Delete("/boards/:id", handlers.DeleteBoard)

	// Column routes
	protected.Post("/columns", handlers.CreateColumn)
	protected.Put("/columns/:id", handlers.UpdateColumn)
	protected.Delete("/columns/:id", handlers.DeleteColumn)

	// Card routes
	protected.Post("/cards", handlers.CreateCard)
	protected.Put("/cards/:id", handlers.UpdateCard)
	protected.Delete("/cards/:id", handlers.DeleteCard)
}
