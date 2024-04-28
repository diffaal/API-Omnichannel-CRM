package entity

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	InteractionId    uint      `json:"interaction_id"`
	SenderId         string    `json:"sender_id"`
	RecipientId      string    `json:"recipient_id"`
	MetaMessageId    string    `json:"mid"`
	Message          string    `json:"message"`
	MessageTimestamp time.Time `json:"message_timestamp"`
	AttachmentType   string    `json:"attachment_type"`
	AttachmentUrl    string    `json:"attachment_url"`
	SentBy           string    `json:"sent_by"`
	IsRead           bool      `json:"is_read"`
	IsDeleted        bool      `json:"is_deleted"`
}
