package handlers

import (
	"tether-server/database"
	"tether-server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateCard - создать новую карточку
func CreateCard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Position    int    `json:"position"`
		Color       string `json:"color"`
		ColumnID    string `json:"column_id"`
		AssigneeID  string `json:"assignee_id"`
		DueDate     string `json:"due_date"`
		// CRM Fields
		LeadName     string  `json:"lead_name"`
		ContactEmail string  `json:"contact_email"`
		ContactPhone string  `json:"contact_phone"`
		Company      string  `json:"company"`
		Value        float64 `json:"value"`
		Priority     string  `json:"priority"`
		Status       string  `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if input.Title == "" || input.ColumnID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Title and column_id are required",
		})
	}

	columnUUID, err := uuid.Parse(input.ColumnID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid column ID",
		})
	}

	// Check column and board access
	var column models.Column
	if err := database.DB.Preload("Board").First(&column, columnUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Column not found",
		})
	}

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

	// Parse assignee ID if provided
	var assigneeID *uuid.UUID
	if input.AssigneeID != "" {
		assigneeUUID, err := uuid.Parse(input.AssigneeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid assignee ID",
			})
		}
		assigneeID = &assigneeUUID
	}

	// Parse due date if provided
	var dueDate *time.Time
	if input.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", input.DueDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid due date format",
			})
		}
		dueDate = &parsedDate
	}

	// Set defaults
	if input.Color == "" {
		input.Color = "#FFFFFF"
	}
	if input.Priority == "" {
		input.Priority = "medium"
	}
	if input.Status == "" {
		input.Status = "new"
	}

	card := models.Card{
		ID:           uuid.New(),
		Title:        input.Title,
		Description:  input.Description,
		Position:     input.Position,
		Color:        input.Color,
		ColumnID:     columnUUID,
		AssigneeID:   assigneeID,
		CreatedByID:  userUUID,
		DueDate:      dueDate,
		LeadName:     input.LeadName,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		Company:      input.Company,
		Value:        input.Value,
		Priority:     input.Priority,
		Status:       input.Status,
		CreatedAt:    time.Now(),
	}

	if err := database.DB.Create(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create card",
		})
	}

	// Load relations
	database.DB.Preload("Assignee").Preload("CreatedBy").First(&card, card.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    card,
	})
}

// UpdateCard - обновить карточку
func UpdateCard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	cardID := c.Params("id")
	cardUUID, err := uuid.Parse(cardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid card ID",
		})
	}

	var card models.Card
	if err := database.DB.Preload("Column.Board").First(&card, cardUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Card not found",
		})
	}

	// Check board access
	board := card.Column.Board
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
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Position    *int    `json:"position"`
		Color       *string `json:"color"`
		ColumnID    *string `json:"column_id"`
		AssigneeID  *string `json:"assignee_id"`
		DueDate     *string `json:"due_date"`
		// CRM Fields
		LeadName     *string  `json:"lead_name"`
		ContactEmail *string  `json:"contact_email"`
		ContactPhone *string  `json:"contact_phone"`
		Company      *string  `json:"company"`
		Value        *float64 `json:"value"`
		Priority     *string  `json:"priority"`
		Status       *string  `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Update fields
	if input.Title != nil {
		card.Title = *input.Title
	}
	if input.Description != nil {
		card.Description = *input.Description
	}
	if input.Position != nil {
		card.Position = *input.Position
	}
	if input.Color != nil {
		card.Color = *input.Color
	}
	if input.ColumnID != nil {
		columnUUID, err := uuid.Parse(*input.ColumnID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid column ID",
			})
		}
		card.ColumnID = columnUUID
	}
	if input.AssigneeID != nil {
		if *input.AssigneeID == "" {
			card.AssigneeID = nil
		} else {
			assigneeUUID, err := uuid.Parse(*input.AssigneeID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"error":   "Invalid assignee ID",
				})
			}
			card.AssigneeID = &assigneeUUID
		}
	}
	if input.DueDate != nil {
		if *input.DueDate == "" {
			card.DueDate = nil
		} else {
			parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", *input.DueDate)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"error":   "Invalid due date format",
				})
			}
			card.DueDate = &parsedDate
		}
	}
	if input.LeadName != nil {
		card.LeadName = *input.LeadName
	}
	if input.ContactEmail != nil {
		card.ContactEmail = *input.ContactEmail
	}
	if input.ContactPhone != nil {
		card.ContactPhone = *input.ContactPhone
	}
	if input.Company != nil {
		card.Company = *input.Company
	}
	if input.Value != nil {
		card.Value = *input.Value
	}
	if input.Priority != nil {
		card.Priority = *input.Priority
	}
	if input.Status != nil {
		card.Status = *input.Status
	}

	card.UpdatedAt = time.Now()

	if err := database.DB.Save(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update card",
		})
	}

	// Load relations
	database.DB.Preload("Assignee").Preload("CreatedBy").First(&card, card.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    card,
	})
}

// DeleteCard - удалить карточку
func DeleteCard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	cardID := c.Params("id")
	cardUUID, err := uuid.Parse(cardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid card ID",
		})
	}

	var card models.Card
	if err := database.DB.Preload("Column.Board").First(&card, cardUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Card not found",
		})
	}

	// Check access - user can delete if they created it or are board owner
	board := card.Column.Board
	if card.CreatedByID != userUUID && board.OwnerID != userUUID {
		if board.WorkspaceID != nil {
			var member models.WorkspaceMember
			if err := database.DB.Where("workspace_id = ? AND user_id = ? AND role IN ?", board.WorkspaceID, userUUID, []string{"owner", "admin"}).First(&member).Error; err != nil {
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

	if err := database.DB.Delete(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete card",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Card deleted successfully",
	})
}
