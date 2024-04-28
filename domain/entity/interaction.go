package entity

import (
	"time"

	"gorm.io/gorm"
)

type Interaction struct {
	gorm.Model
	PlatformId      string    `json:"platform_id"`
	ReporterId      uint      `json:"reporter_id"`
	ConversationId  string    `json:"conversation_id"`
	MentionMediaId  string    `json:"media_id"`
	MentionMediaUrl string    `json:"media_url"`
	AgentId         string    `json:"agent_id"`
	Status          string    `json:"status"`
	Platform        string    `json:"platform"`
	InteractionType string    `json:"interaction_type"`
	Latitude        string    `json:"Latitude"`
	Longitude       string    `json:"Longitude"`
	Duration        time.Time `json:"duration"`
}

type GeotagInformation struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
