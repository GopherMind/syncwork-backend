package chats

import (
	"log"
	"time"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/supabase-community/postgrest-go"
)

// Message - структура одного сообщения с данными отправителя
type Message struct {
	ID        string  `json:"id"`
	ChatID    string  `json:"chat_id"`
	Message   string  `json:"message"`
	SenderID  string  `json:"sender_id"`
	CreatedAt string  `json:"created_at,omitempty"`
	Profiles  Profile `json:"profiles"` // Данные профиля отправителя
}

// Profile - данные профиля отправителя
type Profile struct {
	Name string  `json:"name"`
	Url  *string `json:"url,omitempty"`
}

func GetChatMessages(c *fiber.Ctx) error {
	chatID := c.Params("id")

	if chatID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "BAD_REQUEST",
			"message": "Chat ID is required",
		})
	}

	var messages []Message

	// Получаем сообщения вместе с данными профиля отправителя через JOIN
	// profiles(name, url) - это синтаксис Supabase для получения связанных данных
	_, err := db.SB.From("messages").
		Select("id, chat_id, message, sender_id, created_at, profiles(name, url)", "", false).
		Eq("chat_id", chatID).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&messages)

	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch messages",
		})
	}

	return c.JSON(fiber.Map{
		"messages": messages,
	})
}

func GetChatMessagesWS(c *websocket.Conn) {

	chatID := c.Params("id")
	if chatID == "" {
		c.WriteJSON(fiber.Map{
			"error":   "BAD_REQUEST",
			"message": "Chat ID is required",
		})
		c.Close()
		return
	}

	log.Printf("WebSocket connected for chat: %s", chatID)

	lastCheck := time.Now()

	for {
		time.Sleep(1 * time.Second)

		var newMessages []Message
		// Получаем новые сообщения вместе с данными профиля
		_, err := db.SB.From("messages").
			Select("id, chat_id, message, sender_id, created_at, profiles(name, url)", "", false).
			Eq("chat_id", chatID).
			Gt("created_at", lastCheck.Format(time.RFC3339)).
			Order("created_at", &postgrest.OrderOpts{Ascending: true}).
			ExecuteTo(&newMessages)

		if err != nil {
			log.Printf("Error fetching new messages: %v", err)
			continue
		}

		if len(newMessages) > 0 {
			for _, msg := range newMessages {
				if err := c.WriteJSON(msg); err != nil {
					log.Printf("Error writing message: %v", err)
					return
				}
			}
			lastCheck = time.Now()
		}
	}
}
