package auth

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
	storage_go "github.com/supabase-community/storage-go"
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

	contentType := image.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	uniqueFilename := time.Now().Format("20060102150405") + "_" + image.Filename

	var profiles []profileRow
	_, err = db.SB.From("profiles").
		Select("url", "", false).
		Eq("id", userID.(string)).
		ExecuteTo(&profiles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch existing profile image URL",
		})
	}

	if len(profiles) > 0 && profiles[0].URL != "" {
		parts := strings.SplitN(profiles[0].URL, "/avatar/", 2)
		if len(parts) == 2 && parts[1] != "" {
			decodedKey, err := url.QueryUnescape(parts[1])
			if err != nil {
				decodedKey = parts[1]
			}

			log.Printf("Key to delete (decoded): %s", decodedKey)

			_, err = db.SB.Storage.RemoveFile("avatar", []string{decodedKey})
			if err != nil {
				log.Printf("Error removing existing profile image: %v", err)
			}
		}
	}

	uploadResult, err := db.SB.Storage.UploadFile("avatar", uniqueFilename, fileReader, storage_go.FileOptions{
		ContentType: &contentType,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to upload new profile image",
		})
	}

	cleanKey := strings.TrimPrefix(uploadResult.Key, "avatar/")

	publicURL := db.SB.Storage.GetPublicUrl("avatar", cleanKey)
	imageURL := publicURL.SignedURL

	_, _, err = db.SB.From("profiles").
		Update(map[string]interface{}{"url": imageURL}, "representation", "").
		Eq("id", userID.(string)).
		Execute()
	if err != nil {
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
