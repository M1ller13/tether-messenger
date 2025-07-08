package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"tether-server/database"
	"tether-server/models"
	"tether-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RegisterInput struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	Bio         string `json:"bio"`
}

func Register(c *fiber.Ctx) error {
	var input struct {
		DisplayName string `json:"display_name"`
		Phone       string `json:"phone"`
	}
	if err := c.BodyParser(&input); err != nil || input.Phone == "" || input.DisplayName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Display name and phone required"})
	}
	code := randomCode(6)
	expires := time.Now().Add(5 * time.Minute)
	vc := models.VerificationCode{
		ID:          uuid.New(),
		Phone:       input.Phone,
		Code:        code,
		ExpiresAt:   expires,
		CreatedAt:   time.Now(),
		DisplayName: input.DisplayName,
	}
	database.DB.Where("phone = ?", input.Phone).Delete(&models.VerificationCode{})
	if err := database.DB.Create(&vc).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "DB error"})
	}
	log.Printf("[MOCK SMS] Код для %s: %s (действителен до %s)", input.Phone, code, expires.Format(time.RFC3339))
	return c.JSON(fiber.Map{"success": true})
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Find user in DB
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid credentials",
		})
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":           user.ID,
				"username":     user.Username,
				"display_name": user.DisplayName,
				"bio":          user.Bio,
				"avatar_url":   user.AvatarURL,
				"created_at":   user.CreatedAt,
			},
		},
	})
}

// Search users by username, display_name, bio (partial match, exclude self)
func SearchUsers(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	query := c.Query("query")
	if len(query) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Query parameter is required",
		})
	}

	var users []models.User
	q := "%" + query + "%"
	if err := database.DB.Where(
		"(username ILIKE ? OR display_name ILIKE ? OR bio ILIKE ?) AND id != ?",
		q, q, q, userID,
	).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to search users",
		})
	}

	result := make([]fiber.Map, 0, len(users))
	for _, u := range users {
		result = append(result, fiber.Map{
			"id":           u.ID,
			"username":     u.Username,
			"display_name": u.DisplayName,
			"bio":          u.Bio,
			"avatar_url":   u.AvatarURL,
			"last_seen":    u.LastSeen,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// Получить профиль текущего пользователя
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":           user.ID,
			"display_name": user.DisplayName,
			"username":     user.Username,
			"bio":          user.Bio,
			"avatar_url":   user.AvatarURL,
			"last_seen":    user.LastSeen,
			"created_at":   user.CreatedAt,
		},
	})
}

// Обновить профиль текущего пользователя
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var input struct {
		DisplayName *string `json:"display_name"`
		Bio         *string `json:"bio"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}
	if input.Bio != nil {
		user.Bio = *input.Bio
	}
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update profile",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":           user.ID,
			"display_name": user.DisplayName,
			"username":     user.Username,
			"bio":          user.Bio,
			"avatar_url":   user.AvatarURL,
			"last_seen":    user.LastSeen,
			"created_at":   user.CreatedAt,
		},
	})
}

// Загрузка аватара пользователя
func UploadAvatar(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "No file uploaded",
		})
	}

	// Создать папку uploads, если не существует
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// Сохраняем файл с уникальным именем
	ext := filepath.Ext(file.Filename)
	avatarName := "avatar_" + userID + ext
	avatarPath := filepath.Join(uploadDir, avatarName)
	if err := c.SaveFile(file, avatarPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save file",
		})
	}

	// Обновляем avatar_url в профиле
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	user.AvatarURL = "/uploads/" + avatarName
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update avatar",
		})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"avatar_url": user.AvatarURL,
	})
}

// POST /api/auth/request-code (для входа существующих пользователей)
func RequestCode(c *fiber.Ctx) error {
	var input struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&input); err != nil || input.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Phone required"})
	}
	code := randomCode(6)
	expires := time.Now().Add(5 * time.Minute)
	vc := models.VerificationCode{
		ID:        uuid.New(),
		Phone:     input.Phone,
		Code:      code,
		ExpiresAt: expires,
		CreatedAt: time.Now(),
	}
	database.DB.Where("phone = ?", input.Phone).Delete(&models.VerificationCode{})
	if err := database.DB.Create(&vc).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "DB error"})
	}
	log.Printf("[MOCK SMS] Код для %s: %s (действителен до %s)", input.Phone, code, expires.Format(time.RFC3339))
	return c.JSON(fiber.Map{"success": true})
}

// POST /api/auth/verify-code
func VerifyCode(c *fiber.Ctx) error {
	var input struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&input); err != nil || input.Phone == "" || input.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Phone and code required"})
	}
	var vc models.VerificationCode
	if err := database.DB.Where("phone = ? AND code = ?", input.Phone, input.Code).First(&vc).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "Invalid code"})
	}
	if time.Now().After(vc.ExpiresAt) {
		database.DB.Delete(&vc)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "Code expired"})
	}
	var user models.User
	if err := database.DB.Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		displayName := "Пользователь"
		if vc.DisplayName != "" {
			displayName = vc.DisplayName
		}
		user = models.User{
			ID:          uuid.New(),
			Phone:       input.Phone,
			DisplayName: displayName,
			Username:    fmt.Sprintf("user%s%d", randomCode(4), time.Now().UnixNano()%10000),
			CreatedAt:   time.Now(),
		}
		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("Ошибка создания пользователя: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "DB error"})
		}
	}
	database.DB.Delete(&vc)
	token, err := utils.GenerateToken(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Token error"})
	}
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "phone": user.Phone, "display_name": user.DisplayName, "username": user.Username}}})
}

func randomCode(n int) string {
	digits := "0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, n)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}
