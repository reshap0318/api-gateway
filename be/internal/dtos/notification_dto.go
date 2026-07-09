package dtos

import (
	"time"

	"github.com/reshap0318/api-gateway/internal/models"
)

type NotificationDTO struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Data      string     `json:"data"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func ToNotificationDTO(n *models.Notification) NotificationDTO {
	return NotificationDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Data:      n.Data,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
	}
}

func ToNotificationDTOList(notifications []models.Notification) []NotificationDTO {
	result := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		result[i] = ToNotificationDTO(&n)
	}
	return result
}
