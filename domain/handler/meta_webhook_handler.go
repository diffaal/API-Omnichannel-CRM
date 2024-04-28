package handler

import (
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/presentation"
	"Omnichannel-CRM/package/response"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type MetaWebhookHandler struct {
	metaWebhookService service.IMetaWebhookService
}

func NewMetaWebhookHandler(metaWebhookService service.IMetaWebhookService) *MetaWebhookHandler {
	metaWebhookHandler := MetaWebhookHandler{
		metaWebhookService: metaWebhookService,
	}
	return &metaWebhookHandler
}

func (mwh *MetaWebhookHandler) ValidateVerificationRequest(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	challengeInt, err := strconv.Atoi(challenge)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if mode != "" && token != "" && challenge != "" {
		if mode == "subscribe" && token == viper.GetString("Webhook.Token") {
			logger.Info("[INFO] WEBHOOK_VERIFIED")
			c.JSON(http.StatusOK, challengeInt)
			return
		} else {
			c.JSON(http.StatusForbidden, nil)
			return
		}
	}
	c.JSON(http.StatusBadRequest, nil)
}

func (mwh *MetaWebhookHandler) FacebookWebhookInteractionHandler(c *gin.Context) {
	var fwr presentation.FacebookWebhookRequest
	errorMessage := make(map[string]string)

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.ResponseInternalServerError(c, nil, errorMessage)
	}
	bodyString := string(bodyBytes)
	logger.Info(bodyString)

	defer c.Request.Body.Close()

	err = json.Unmarshal(bodyBytes, &fwr)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Webhook Message Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, resMessage, err := mwh.metaWebhookService.FacebookInteractionService(&fwr)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Webhook Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	err = mwh.metaWebhookService.WebsocketSendService(resMessage)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Whatsapp Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (mwh *MetaWebhookHandler) InstagramWebhookInteractionHandler(c *gin.Context) {
	var iwr presentation.InstagramWebhookRequest
	errorMessage := make(map[string]string)

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.ResponseInternalServerError(c, nil, errorMessage)
	}
	bodyString := string(bodyBytes)
	logger.Info(bodyString)

	defer c.Request.Body.Close()

	err = json.Unmarshal(bodyBytes, &iwr)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Webhook Message Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, resMessage, err := mwh.metaWebhookService.InstagramInteractionService(&iwr)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Webhook Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	err = mwh.metaWebhookService.WebsocketSendService(resMessage)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Whatsapp Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (mwh *MetaWebhookHandler) WhatsappMessageInteractionHandler(c *gin.Context) {
	var wir presentation.WhatsappInteractionRequest
	errorMessage := make(map[string]string)

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.ResponseInternalServerError(c, nil, errorMessage)
	}
	bodyString := string(bodyBytes)
	logger.Info(bodyString)

	defer c.Request.Body.Close()

	err = json.Unmarshal(bodyBytes, &wir)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Whatsapp Message Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, resMessage, err := mwh.metaWebhookService.WhatsappMessageInteractionService(&wir)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Whatsapp Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	err = mwh.metaWebhookService.WebsocketSendService(resMessage)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Whatsapp Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}
