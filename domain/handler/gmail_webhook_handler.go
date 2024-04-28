package handler

import (
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/presentation"
	"Omnichannel-CRM/package/response"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GmailHandler struct {
	emailService service.IEmailService
	historyId    int
}

func NewGmailHandler(emailService service.IEmailService, historyId int) GmailHandler {
	return GmailHandler{
		emailService: emailService,
		historyId:    historyId,
	}
}

type AttachmentMetaData struct {
	Name          string `json:"name"`
	ContentType   string `json:"contentType"`
	ContentLength int    `json:"contentLength"`
	ID            string `json:"id"`
}

type WebhookPayload struct {
	MessageID           string               `json:"messageId"`
	WebhookID           string               `json:"webhookId"`
	EventName           string               `json:"eventName"`
	WebhookName         string               `json:"webhookName"`
	InboxID             string               `json:"inboxId"`
	DomainID            interface{}          `json:"domainId"` // Set as interface{} as it can be null
	EmailID             string               `json:"emailId"`
	CreatedAt           time.Time            `json:"createdAt"`
	To                  []string             `json:"to"`
	From                string               `json:"from"`
	CC                  []string             `json:"cc"`
	BCC                 []string             `json:"bcc"`
	Subject             string               `json:"subject"`
	AttachmentMetaDatas []AttachmentMetaData `json:"attachmentMetaDatas"`
}

func (handler *GmailHandler) SetHistoryId(historyId int) {
	handler.historyId = historyId
}

func (handler *GmailHandler) Webhook(c *gin.Context) {
	var req presentation.GmailInteractionRequest
	errorMessage := make(map[string]string)

	err := c.BindJSON(&req)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Gmail Message Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	logger.Info(fmt.Sprintf("====== Received Gmail webhook payload: %+v", req))

	historyId, err := handler.emailService.ProcessWebhook(req.Message.Data, handler.historyId)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Gmail Message Interaction] Internal Error: %+v", err))
		// response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	handler.SetHistoryId(int(historyId))
	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})
}
