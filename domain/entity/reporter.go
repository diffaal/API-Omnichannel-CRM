package entity

import "gorm.io/gorm"

type Reporter struct {
	gorm.Model
	MetaReporterId   string `json:"meta_reporter_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Gender           string `json:"gender"`
	Address          string `json:"address"`
	ReporterStatus   string `json:"reporter_status"`
	PlatformUsername string `json:"platform_username"`
	IsDeleted        bool   `json:"is_deleted"`
}
