package handlers

import (
	"tether-server/database"
	"tether-server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/chats
func GetChats(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid user ID"})
	}
	var chats []models.Chat
	if err := database.DB.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Failed to get chats"})
	}
	return c.JSON(fiber.Map{"success": true, "data": chats})
}

// POST /api/chats
func CreateChat(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid user ID"})
	}
	var input struct {
		OtherUserID uuid.UUID `json:"other_user_id"`
	}
	if err := c.BodyParser(&input); err != nil || input.OtherUserID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}
	// Проверка на существование чата
	var chat models.Chat
	if err := database.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, input.OtherUserID, input.OtherUserID, userID,
	).First(&chat).Error; err == nil {
		return c.JSON(fiber.Map{"success": true, "data": chat})
	}
	// Создать новый чат
	newChat := models.Chat{
		ID:        uuid.New(),
		User1ID:   userID,
		User2ID:   input.OtherUserID,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&newChat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Failed to create chat"})
	}
	return c.JSON(fiber.Map{"success": true, "data": newChat})
}

// GET /api/chats/:chatId/messages
func GetMessages(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid user ID"})
	}
	chatID, err := uuid.Parse(c.Params("chatId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid chat ID"})
	}
	// Проверка доступа
	var chat models.Chat
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?)", chatID, userID, userID).First(&chat).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Access denied"})
	}
	var messages []models.Message
	if err := database.DB.Where("chat_id = ?", chatID).Order("created_at asc").Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Failed to get messages"})
	}
	return c.JSON(fiber.Map{"success": true, "data": messages})
}

// POST /api/messages
func SendMessage(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid user ID"})
	}
	var input struct {
		ChatID  uuid.UUID `json:"chat_id"`
		Content string    `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil || input.ChatID == uuid.Nil || input.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}
	// Проверка доступа
	var chat models.Chat
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?)", input.ChatID, userID, userID).First(&chat).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Access denied"})
	}
	msg := models.Message{
		ID:        uuid.New(),
		ChatID:    input.ChatID,
		SenderID:  userID,
		Content:   input.Content,
		CreatedAt: time.Now(),
		IsRead:    false,
	}
	if err := database.DB.Create(&msg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Failed to send message"})
	}
	return c.JSON(fiber.Map{"success": true, "data": msg})
}
