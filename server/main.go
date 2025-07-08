package main

import (
	"log"
	"tether-server/config"
	"tether-server/database"
	"tether-server/routes"
	"tether-server/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	// Подключаемся к базе данных
	database.ConnectDB()

	// Создаем Fiber приложение
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Настраиваем маршруты
	routes.SetupRoutes(app)

	// WebSocket маршрут
	app.Get("/ws", ws.WebSocketHandler())

	// Запускаем WebSocket hub в горутине
	go ws.Run()

	// Запускаем сервер
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	log.Fatal(app.Listen(":" + config.AppConfig.ServerPort))
}
