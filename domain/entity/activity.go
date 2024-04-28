package entity

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	AgentId string `json:"agent_id"`
	ActivityStatus string `json:"activity_status"`
}
