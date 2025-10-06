package handlers

import (
	"tether-server/database"
	"tether-server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateBoard - создать новую доску
func CreateBoard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
		WorkspaceID string `json:"workspace_id"`
		IsPublic    bool   `json:"is_public"`
		Color       string `json:"color"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Board name is required",
		})
	}

	// Validate board type
	if input.Type == "" {
		input.Type = "personal"
	}
	if input.Type != "personal" && input.Type != "team" && input.Type != "crm" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid board type",
		})
	}

	// Validate workspace access if provided
	var workspaceID *uuid.UUID
	if input.WorkspaceID != "" {
		wsUUID, err := uuid.Parse(input.WorkspaceID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid workspace ID",
			})
		}

		// Check if user has access to workspace
		var member models.WorkspaceMember
		if err := database.DB.Where("workspace_id = ? AND user_id = ?", wsUUID, userUUID).First(&member).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Access denied to workspace",
			})
		}

		workspaceID = &wsUUID
	}

	// Set default color
	if input.Color == "" {
		input.Color = "#3B82F6"
	}

	board := models.Board{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		OwnerID:     userUUID,
		WorkspaceID: workspaceID,
		IsPublic:    input.IsPublic,
		Color:       input.Color,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&board).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create board",
		})
	}

	// Create default columns for the board
	defaultColumns := []string{"To Do", "In Progress", "Done"}
	for i, columnName := range defaultColumns {
		column := models.Column{
			ID:        uuid.New(),
			Name:      columnName,
			Position:  i,
			Color:     "#6B7280",
			BoardID:   board.ID,
			CreatedAt: time.Now(),
		}
		database.DB.Create(&column)
	}

	// Load relations
	database.DB.Preload("Owner").Preload("Workspace").Preload("Columns").First(&board, board.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    board,
	})
}

// GetBoards - получить доски пользователя
func GetBoards(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	// Get personal boards
	var personalBoards []models.Board
	database.DB.Where("owner_id = ? AND workspace_id IS NULL", userUUID).
		Preload("Owner").Preload("Columns").
		Find(&personalBoards)

	// Get workspace boards where user is a member
	var workspaceBoards []models.Board
	database.DB.Joins("JOIN workspace_members ON boards.workspace_id = workspace_members.workspace_id").
		Where("workspace_members.user_id = ?", userUUID).
		Preload("Owner").Preload("Workspace").Preload("Columns").
		Find(&workspaceBoards)

	// Combine boards
	allBoards := append(personalBoards, workspaceBoards...)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    allBoards,
	})
}

// GetBoard - получить конкретную доску
func GetBoard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	boardID := c.Params("id")
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid board ID",
		})
	}

	var board models.Board
	if err := database.DB.Preload("Owner").Preload("Workspace").Preload("Columns.Cards.Assignee").Preload("Columns.Cards.CreatedBy").First(&board, boardUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Board not found",
		})
	}

	// Check access
	if board.OwnerID != userUUID && !board.IsPublic {
		// Check if user is member of workspace
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

	return c.JSON(fiber.Map{
		"success": true,
		"data":    board,
	})
}

// UpdateBoard - обновить доску
func UpdateBoard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	boardID := c.Params("id")
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid board ID",
		})
	}

	var board models.Board
	if err := database.DB.First(&board, boardUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Board not found",
		})
	}

	// Check ownership
	if board.OwnerID != userUUID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Only board owner can update",
		})
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		IsPublic    *bool   `json:"is_public"`
		Color       *string `json:"color"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Update fields
	if input.Name != nil {
		board.Name = *input.Name
	}
	if input.Description != nil {
		board.Description = *input.Description
	}
	if input.IsPublic != nil {
		board.IsPublic = *input.IsPublic
	}
	if input.Color != nil {
		board.Color = *input.Color
	}

	board.UpdatedAt = time.Now()

	if err := database.DB.Save(&board).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update board",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    board,
	})
}

// DeleteBoard - удалить доску
func DeleteBoard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	boardID := c.Params("id")
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid board ID",
		})
	}

	var board models.Board
	if err := database.DB.First(&board, boardUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Board not found",
		})
	}

	// Check ownership
	if board.OwnerID != userUUID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Only board owner can delete",
		})
	}

	if err := database.DB.Delete(&board).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete board",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Board deleted successfully",
	})
}
