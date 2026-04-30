package auth

import (
	"log"
	"net/url"
	"time"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
)

type profileRow struct {
	URL string `json:"url"`
}

func LoadedImage(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error":   "FORBIDDEN",
			"message": "Failed to retrieve the uploaded file",
		})
	}

	fileReader, err := image.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to open uploaded file",
		})
	}
	defer fileReader.Close()

	uniqueFilename := time.Now().Format("20060102150405") + "_" + image.Filename

	var profiles []profileRow
	_, err = db.SB.From("profiles").
		Select("url", "", false).
		Eq("id", userID.(string)).
		ExecuteTo(&profiles)
	if err != nil {
		log.Printf("Error fetching existing profile image URL: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch existing profile image URL",
		})
	}

	if len(profiles) > 0 && profiles[0].URL != "" {
		existingPath := url.PathEscape(profiles[0].URL)
		_, err := db.SB.Storage.RemoveFile("avatar", []string{existingPath})
		if err != nil {
			log.Printf("Error removing existing profile image: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   "INTERNAL_ERROR",
				"message": "Failed to remove existing profile image",
			})
		}
	}

	uploadResult, err := db.SB.Storage.UploadFile("avatar", uniqueFilename, fileReader)
	if err != nil {
		log.Printf("Error uploading new profile image: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to upload new profile image",
		})
	}
	imageURL := uploadResult.Key

	_, _, err = db.SB.From("profiles").
		Update(map[string]interface{}{"url": imageURL}, "representation", "").
		Eq("id", userID.(string)).
		Execute()
	if err != nil {
		log.Printf("Error updating profile with new image URL: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to update profile with new image URL",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Profile image updated successfully",
		"imageURL": imageURL,
	})
}
