package presentation

type UpdateReporterRequest struct {
	ReporterId       uint   `json:"reporter_id" binding:"required"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Gender           string `json:"gender"`
	Address          string `json:"address"`
	PlatformUsername string `json:"platform_username"`
}
