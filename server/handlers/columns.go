package handlers

import (
	"tether-server/database"
	"tether-server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateColumn - создать новую колонку
func CreateColumn(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	var input struct {
		Name     string `json:"name"`
		Position int    `json:"position"`
		Color    string `json:"color"`
		BoardID  string `json:"board_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if input.Name == "" || input.BoardID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Name and board_id are required",
		})
	}

	boardUUID, err := uuid.Parse(input.BoardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid board ID",
		})
	}

	// Check board access
	var board models.Board
	if err := database.DB.First(&board, boardUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Board not found",
		})
	}

	// Check if user has access to board
	if board.OwnerID != userUUID && !board.IsPublic {
		if board.WorkspaceID != nil {
			var member models.WorkspaceMember
			if err := database.DB.Where("workspace_id = ? AND user_id = ?", board.WorkspaceID, userUUID).First(&member).Error; err != nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"error":   "Access denied",
				})
			}
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Access denied",
			})
		}
	}

	// Set default color
	if input.Color == "" {
		input.Color = "#6B7280"
	}

	column := models.Column{
		ID:        uuid.New(),
		Name:      input.Name,
		Position:  input.Position,
		Color:     input.Color,
		BoardID:   boardUUID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&column).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create column",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    column,
	})
}

// UpdateColumn - обновить колонку
func UpdateColumn(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	columnID := c.Params("id")
	columnUUID, err := uuid.Parse(columnID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid column ID",
		})
	}

	var column models.Column
	if err := database.DB.Preload("Board").First(&column, columnUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Column not found",
		})
	}

	// Check board access
	board := column.Board
	if board.OwnerID != userUUID && !board.IsPublic {
		if board.WorkspaceID != nil {
			var member models.WorkspaceMember
			if err := database.DB.Where("workspace_id = ? AND user_id = ?", board.WorkspaceID, userUUID).First(&member).Error; err != nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"error":   "Access denied",
				})
			}
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Access denied",
			})
		}
	}

	var input struct {
		Name     *string `json:"name"`
		Position *int    `json:"position"`
		Color    *string `json:"color"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Update fields
	if input.Name != nil {
		column.Name = *input.Name
	}
	if input.Position != nil {
		column.Position = *input.Position
	}
	if input.Color != nil {
		column.Color = *input.Color
	}

	column.UpdatedAt = time.Now()

	if err := database.DB.Save(&column).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update column",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    column,
	})
}

// DeleteColumn - удалить колонку
func DeleteColumn(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	columnID := c.Params("id")
	columnUUID, err := uuid.Parse(columnID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid column ID",
		})
	}

	var column models.Column
	if err := database.DB.Preload("Board").First(&column, columnUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Column not found",
		})
	}

	// Check board ownership
	board := column.Board
	if board.OwnerID != userUUID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Only board owner can delete columns",
		})
	}

	if err := database.DB.Delete(&column).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete column",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Column deleted successfully",
	})
}
