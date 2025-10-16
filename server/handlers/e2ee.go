package handlers

import (
	"net/http"
	"time"

	"tether-server/database"
	"tether-server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PublishDeviceKeys allows an authenticated user to publish device key bundle (identity, signed prekey, prekey signature) and an optional batch of one-time prekeys.
// Private keys must NEVER be sent here — only public materials.
func PublishDeviceKeys(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "invalid user id"})
	}

	var input struct {
		DeviceID              string `json:"device_id"`
		IdentityKeyPublic     string `json:"identity_key_public"`
		SignedPreKeyPublic    string `json:"signed_prekey_public"`
		SignedPreKeySignature string `json:"signed_prekey_signature"`
		OneTimePreKeys        []struct {
			KeyID     int    `json:"key_id"`
			PublicKey string `json:"public_key"`
		} `json:"one_time_prekeys"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	if input.DeviceID == "" || input.IdentityKeyPublic == "" || input.SignedPreKeyPublic == "" || input.SignedPreKeySignature == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "missing required fields"})
	}

	// Upsert DeviceKey (by user+device)
	var device models.DeviceKey
	err = database.DB.Where("user_id = ? AND device_id = ?", userID, input.DeviceID).First(&device).Error
	now := time.Now()
	if err != nil {
		device = models.DeviceKey{
			ID:                    uuid.New(),
			UserID:                userID,
			DeviceID:              input.DeviceID,
			IdentityKeyPublic:     input.IdentityKeyPublic,
			SignedPreKeyPublic:    input.SignedPreKeyPublic,
			SignedPreKeySignature: input.SignedPreKeySignature,
			Active:                true,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		if err := database.DB.Create(&device).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "failed to save device keys"})
		}
	} else {
		device.IdentityKeyPublic = input.IdentityKeyPublic
		device.SignedPreKeyPublic = input.SignedPreKeyPublic
		device.SignedPreKeySignature = input.SignedPreKeySignature
		device.Active = true
		device.UpdatedAt = now
		if err := database.DB.Save(&device).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "failed to update device keys"})
		}
	}

	// Replace existing unused OTPKs if provided: simple approach — insert new ones, ignore conflicts by (device_id,key_id)
	if len(input.OneTimePreKeys) > 0 {
		for _, k := range input.OneTimePreKeys {
			if k.KeyID <= 0 || k.PublicKey == "" {
				continue
			}
			otp := models.OneTimePreKey{
				ID:          uuid.New(),
				DeviceKeyID: device.ID,
				KeyID:       k.KeyID,
				PublicKey:   k.PublicKey,
				Used:        false,
				CreatedAt:   now,
			}
			// Try create; if unique conflict, skip
			_ = database.DB.Where("device_key_id = ? AND key_id = ?", device.ID, k.KeyID).First(&models.OneTimePreKey{}).Error
			_ = database.DB.Clauses().Create(&otp)
		}
	}

	return c.JSON(fiber.Map{"success": true})
}

// FetchPreKeyBundle returns public bundle for initiating session with a target user (optionally specific device).
// Includes: identity key, signed prekey (+signature), and one available one-time prekey (and marks it used).
func FetchPreKeyBundle(c *fiber.Ctx) error {
	targetUserIDStr := c.Params("userId")
	if targetUserIDStr == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "userId required"})
	}
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid userId"})
	}

	deviceID := c.Query("device_id")

	var device models.DeviceKey
	q := database.DB.Where("user_id = ? AND active = ?", targetUserID, true)
	if deviceID != "" {
		q = q.Where("device_id = ?", deviceID)
	}
	if err := q.First(&device).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"success": false, "error": "no device keys"})
	}

	// Get one unused OTPK if available
	var otp models.OneTimePreKey
	if err := database.DB.Where("device_key_id = ? AND used = ?", device.ID, false).Order("key_id asc").First(&otp).Error; err == nil {
		// mark used
		now := time.Now()
		otp.Used = true
		otp.UsedAt = &now
		_ = database.DB.Save(&otp).Error
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":                 device.UserID,
			"device_id":               device.DeviceID,
			"identity_key_public":     device.IdentityKeyPublic,
			"signed_prekey_public":    device.SignedPreKeyPublic,
			"signed_prekey_signature": device.SignedPreKeySignature,
			"one_time_prekey": func() interface{} {
				if otp.ID != uuid.Nil {
					return fiber.Map{"key_id": otp.KeyID, "public_key": otp.PublicKey}
				}
				return nil
			}(),
		},
	})
}
