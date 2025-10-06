package handlers

import (
	"fmt"
	"log"
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
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
}

// Register - email-based registration
func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate input
	if input.Email == "" || input.Password == "" || input.DisplayName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email, password and display name are required",
		})
	}

	if !utils.IsValidEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid email format",
		})
	}

	if len(input.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password must be at least 6 characters",
		})
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   "User with this email already exists",
		})
	}

	// Generate username if not provided
	username := input.Username
	if username == "" {
		username = fmt.Sprintf("user%d", time.Now().UnixNano()%1000000)
	}

	// Check username uniqueness
	var usernameUser models.User
	if err := database.DB.Where("username = ?", username).First(&usernameUser).Error; err == nil {
		username = fmt.Sprintf("%s%d", username, time.Now().UnixNano()%1000)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to process password",
		})
	}

	// Create user
	user := models.User{
		ID:          uuid.New(),
		Email:       input.Email,
		Password:    hashedPassword,
		Username:    username,
		DisplayName: input.DisplayName,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create user",
		})
	}

	// Generate email verification token
	token, err := utils.GenerateEmailToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate verification token",
		})
	}

	// Create email verification record
	emailVerification := models.EmailVerification{
		ID:        uuid.New(),
		Email:     input.Email,
		Token:     token,
		Type:      "signup",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&emailVerification).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create verification record",
		})
	}

	// Send verification email
	if err := utils.SendVerificationEmail(input.Email, token, input.DisplayName); err != nil {
		log.Printf("Failed to send verification email: %v", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

// Login - email-based login
func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid credentials",
		})
	}

	// Check password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid credentials",
		})
	}

	// Check if email is verified
	if !user.EmailVerified {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Please verify your email before logging in",
		})
	}

	// Generate token pair
	tokenPair, err := utils.GenerateTokenPair(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate tokens",
		})
	}

	// Store refresh token in database
	refreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&refreshToken).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to store refresh token",
		})
	}

	// Update last seen
	now := time.Now()
	user.LastSeen = &now
	database.DB.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"user": fiber.Map{
				"id":             user.ID,
				"email":          user.Email,
				"username":       user.Username,
				"display_name":   user.DisplayName,
				"bio":            user.Bio,
				"avatar_url":     user.AvatarURL,
				"email_verified": user.EmailVerified,
				"created_at":     user.CreatedAt,
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

// VerifyEmail - verify email address with token
func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Token is required",
		})
	}

	// Find verification record
	var emailVerification models.EmailVerification
	if err := database.DB.Where("token = ? AND type = ? AND used = ?", token, "signup", false).First(&emailVerification).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid or expired verification token",
		})
	}

	// Check if token is expired
	if time.Now().After(emailVerification.ExpiresAt) {
		database.DB.Delete(&emailVerification)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Verification token has expired",
		})
	}

	// Find user and mark email as verified
	var user models.User
	if err := database.DB.Where("email = ?", emailVerification.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	user.EmailVerified = true
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to verify email",
		})
	}

	// Mark verification as used
	emailVerification.Used = true
	database.DB.Save(&emailVerification)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email verified successfully",
	})
}

// RequestPasswordReset - request password reset
func RequestPasswordReset(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if !utils.IsValidEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid email format",
		})
	}

	// Check if user exists
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Don't reveal if user exists or not for security
		return c.JSON(fiber.Map{
			"success": true,
			"message": "If the email exists, a password reset link has been sent",
		})
	}

	// Generate reset token
	token, err := utils.GenerateEmailToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate reset token",
		})
	}

	// Create password reset record
	passwordReset := models.EmailVerification{
		ID:        uuid.New(),
		Email:     input.Email,
		Token:     token,
		Type:      "password_reset",
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		CreatedAt: time.Now(),
	}

	// Delete any existing reset tokens for this email
	database.DB.Where("email = ? AND type = ?", input.Email, "password_reset").Delete(&models.EmailVerification{})

	if err := database.DB.Create(&passwordReset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create reset record",
		})
	}

	// Send reset email
	if err := utils.SendPasswordResetEmail(input.Email, token, user.DisplayName); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword - reset password with token
func ResetPassword(c *fiber.Ctx) error {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if len(input.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password must be at least 6 characters",
		})
	}

	// Find reset record
	var passwordReset models.EmailVerification
	if err := database.DB.Where("token = ? AND type = ? AND used = ?", input.Token, "password_reset", false).First(&passwordReset).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid or expired reset token",
		})
	}

	// Check if token is expired
	if time.Now().After(passwordReset.ExpiresAt) {
		database.DB.Delete(&passwordReset)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Reset token has expired",
		})
	}

	// Find user and update password
	var user models.User
	if err := database.DB.Where("email = ?", passwordReset.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to process password",
		})
	}

	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update password",
		})
	}

	// Mark reset token as used
	passwordReset.Used = true
	database.DB.Save(&passwordReset)

	// Revoke all refresh tokens for this user
	database.DB.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Update("revoked", true)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password reset successfully",
	})
}

// RefreshToken - refresh access token
func RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Find refresh token
	var refreshToken models.RefreshToken
	if err := database.DB.Where("token = ? AND revoked = ?", input.RefreshToken, false).First(&refreshToken).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid refresh token",
		})
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		refreshToken.Revoked = true
		database.DB.Save(&refreshToken)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Refresh token has expired",
		})
	}

	// Generate new token pair
	tokenPair, err := utils.GenerateTokenPair(refreshToken.UserID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate tokens",
		})
	}

	// Revoke old refresh token
	refreshToken.Revoked = true
	database.DB.Save(&refreshToken)

	// Store new refresh token
	newRefreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshToken.UserID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&newRefreshToken).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to store refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
		},
	})
}

// Logout - revoke refresh token
func Logout(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Revoke refresh token
	database.DB.Model(&models.RefreshToken{}).Where("token = ?", input.RefreshToken).Update("revoked", true)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}
