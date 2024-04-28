package handler

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/presentation"
	"Omnichannel-CRM/package/response"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	interactionService service.IInteractionService
}

func NewInteractionHandler(interactionService service.IInteractionService) *InteractionHandler {
	interactionHandler := InteractionHandler{
		interactionService: interactionService,
	}
	return &interactionHandler
}

func (ih *InteractionHandler) GetDashboardInteractionList(c *gin.Context) {
	var channelAccount entity.ChannelAccount
	errorMessage := make(map[string]string)

	channelAccountData := c.Keys["channel_account"]
	channelAccountJson, err := json.Marshal(channelAccountData)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Interaction List] Invalid Channel Account Data from Token: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	err = json.Unmarshal(channelAccountJson, &channelAccount)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Interaction List] Invalid Channel Account Data from Token: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	filters, err := presentation.ParseGetListInteractionFilters(c)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Interaction List] Invalid Query Params: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	result, err := ih.interactionService.GetInteractionList(filters, &channelAccount)

	if result["errorStatus"] != nil {
		errorMessage["errorMessage"] = result["errorMessage"].(string)
		errorMessage["errorStatus"] = result["errorStatus"].(string)

		if errorMessage["errorStatus"] == enum.DATA_NOT_FOUND_STATUS {
			response.ResponseNotFound(c, nil, errorMessage)
			return
		} else {
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) ClaimInteractionByAgent(c *gin.Context) {
	result := make(map[string]interface{})
	userId := c.GetString("user_id")
	var cir presentation.ClaimInteractionRequest
	errorMessage := make(map[string]string)

	err := c.BindJSON(&cir)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Claim Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	validation := cir.ValidatePayload()
	if validation["errorStatus"] != "" {
		logger.Info("[FAILED][Claim Interaction] Invalid Payload")
		response.ResponseInvalidRequest(c, nil, validation)
		return
	}

	interaction, err := ih.interactionService.UpdateInteractionStatusByAgent(&cir, userId, enum.IN_PROGRESS)
	if errors.Is(err, enum.ERROR_DATA_NOT_FOUND) {
		errorMessage["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		errorMessage["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		response.ResponseNotFound(c, nil, errorMessage)
		return

	} else if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	result["interaction_id"] = interaction.ID
	result["agent_id"] = interaction.AgentId
	result["interaction_status"] = interaction.Status

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) MessengerSendMessage(c *gin.Context) {
	var msmr presentation.MetaSendMessageRequest
	result := make(map[string]interface{})
	errorMessage := make(map[string]string)
	var message *entity.Message
	var channelAccount entity.ChannelAccount

	channelAccountData := c.Keys["channel_account"]
	channelAccountJson, err := json.Marshal(channelAccountData)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Interaction List] Invalid Channel Account Data from Token: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	err = json.Unmarshal(channelAccountJson, &channelAccount)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Interaction List] Invalid Channel Account Data from Token: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	err = c.BindJSON(&msmr)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Claim Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	platform := msmr.Platform
	if platform == enum.FACEBOOK || platform == enum.IG {
		result, message, err = ih.interactionService.MessengerSendMessagetoMeta(&msmr, &channelAccount)
		if errors.Is(err, enum.USER_DO_NOT_HAVE_CHANNEL_ACCOUNT) {
			errorMessage["errorMessage"] = enum.USER_DO_NOT_HAVE_CHANNEL_ACCOUNT_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseBadRequest(c, nil, errorMessage)
			return

		} else if errors.Is(err, enum.PLATFORM_ID_NOT_SET) {
			errorMessage["errorMessage"] = enum.PLATFORM_ID_NOT_SET_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return

		} else if errors.Is(err, enum.PLATFORM_ACCESS_TOKEN_NOT_SET) {
			errorMessage["errorMessage"] = enum.PLATFORM_ACCESS_TOKEN_NOT_SET_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return

		} else if errors.Is(err, enum.CHANNEL_ACCOUNT_NOT_MATCH) {
			errorMessage["errorMessage"] = enum.CHANNEL_ACCOUNT_NOT_MATCH_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseBadRequest(c, nil, errorMessage)
			return

		} else if err != nil {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta] Internal Error: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}

	} else if platform == enum.WA {
		result, message, err = ih.interactionService.WhatsappSendMessagetoMeta(&msmr, &channelAccount)
		if errors.Is(err, enum.USER_DO_NOT_HAVE_CHANNEL_ACCOUNT) {
			errorMessage["errorMessage"] = enum.USER_DO_NOT_HAVE_CHANNEL_ACCOUNT_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseBadRequest(c, nil, errorMessage)
			return

		} else if errors.Is(err, enum.PLATFORM_ID_NOT_SET) {
			errorMessage["errorMessage"] = enum.PLATFORM_ID_NOT_SET_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return

		} else if errors.Is(err, enum.PLATFORM_ACCESS_TOKEN_NOT_SET) {
			errorMessage["errorMessage"] = enum.PLATFORM_ACCESS_TOKEN_NOT_SET_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return

		} else if errors.Is(err, enum.CHANNEL_ACCOUNT_NOT_MATCH) {
			errorMessage["errorMessage"] = enum.CHANNEL_ACCOUNT_NOT_MATCH_MSG
			errorMessage["errorStatus"] = enum.FAILED_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta]: %+v", err))
			response.ResponseBadRequest(c, nil, errorMessage)
			return

		} else if err != nil {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Messenger Send Message to Meta] Internal Error: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}

	} else if platform == enum.LIVE_CHAT {
		result, message, err = ih.interactionService.LiveChatSendMessage(&msmr)
		if err != nil {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Whatsapp Send Message to Meta] Internal Error: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}

	} else if platform == enum.EMAIL {
		result, _, err = ih.interactionService.SendMessageToEmail(presentation.MessengerSendEmailRequest{
			InteractionId: msmr.InteractionId,
			Message:       msmr.Message,
		})
	}

	if platform != enum.EMAIL {
		err = ih.interactionService.WebsocketSendService(*message)
		if err != nil {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Whatsapp Message Interaction] Internal Error: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) LiveChatSendMessage(c *gin.Context) {
	var msmr presentation.MetaSendMessageRequest
	result := make(map[string]interface{})
	errorMessage := make(map[string]string)
	var message *entity.Message

	err := c.BindJSON(&msmr)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Live Chat Send Message] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	platform := msmr.Platform

	if platform == enum.LIVE_CHAT {
		result, message, err = ih.interactionService.LiveChatSendMessage(&msmr)
		if err != nil {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Live Chat Send Message] Internal Error: %+v", err))
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}

	} else {
		errorMessage["errorMessage"] = enum.INVALID_PLATFORM_MSG
		errorMessage["errorStatus"] = enum.FAILED_STATUS
		logger.Info("[FAILED][Live Chat Send Message] Invalid Platform")
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	err = ih.interactionService.WebsocketSendService(*message)
	if err != nil {
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Live Chat Message Interaction] Internal Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) GetInteractionMessages(c *gin.Context) {
	errorMessage := make(map[string]string)

	interactionIdQuery := c.Query("interaction_id")
	if interactionIdQuery == "" {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Claim Interaction] Invalid Query. Required interaction_id"))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	interactionId, err := strconv.ParseUint(interactionIdQuery, 10, 64)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Claim Interaction] Invalid Value of Query interaction_id"))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}
	interactionIdUint := uint(interactionId)

	result, err := ih.interactionService.GetInteractionMessages(interactionIdUint)
	if result["errorStatus"] != nil {
		errorMessage["errorMessage"] = result["errorMessage"].(string)
		errorMessage["errorStatus"] = result["errorStatus"].(string)

		if errorMessage["errorStatus"] == enum.DATA_NOT_FOUND_STATUS {
			response.ResponseNotFound(c, nil, errorMessage)
			return
		} else {
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) GetAgentInteractions(c *gin.Context) {
	userId := c.GetString("user_id")
	errorMessage := make(map[string]string)

	filters, err := presentation.ParseGetListInteractionFilters(c)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Agent Interactions] Invalid Query Params: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	result, _ := ih.interactionService.GetAgentInteractions(userId, filters)
	if result["errorStatus"] != nil {
		errorMessage["errorMessage"] = result["errorMessage"].(string)
		errorMessage["errorStatus"] = result["errorStatus"].(string)

		if errorMessage["errorStatus"] == enum.DATA_NOT_FOUND_STATUS {
			response.ResponseNotFound(c, nil, errorMessage)
			return
		} else {
			response.ResponseInternalServerError(c, nil, errorMessage)
			return
		}
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) CloseInteractionByAgent(c *gin.Context) {
	result := make(map[string]interface{})
	userId := c.GetString("user_id")
	var cir presentation.ClaimInteractionRequest
	errorMessage := make(map[string]string)

	err := c.BindJSON(&cir)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Claim Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	validation := cir.ValidatePayload()
	if validation["errorStatus"] != "" {
		logger.Info("[FAILED][Claim Interaction] Invalid Payload")
		response.ResponseInvalidRequest(c, nil, validation)
		return
	}

	interaction, err := ih.interactionService.UpdateInteractionStatusByAgent(&cir, userId, enum.CLOSED)
	if errors.Is(err, enum.ERROR_DATA_NOT_FOUND) {
		errorMessage["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		errorMessage["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		response.ResponseNotFound(c, nil, errorMessage)
		return

	} else if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	err = ih.interactionService.SendClosedInteractionData(interaction)
	if err != nil {
		ih.interactionService.UpdateInteractionStatusByAgent(&cir, userId, enum.IN_PROGRESS)
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	result["interaction_id"] = interaction.ID
	result["agent_id"] = interaction.AgentId
	result["interaction_status"] = interaction.Status

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) GetClosedInteractionsData(c *gin.Context) {
	errorMessage := make(map[string]string)

	result, err := ih.interactionService.GetClosedInteractionsData()
	if errors.Is(err, enum.ERROR_DATA_NOT_FOUND) {
		errorMessage["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		errorMessage["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		response.ResponseNotFound(c, nil, errorMessage)
		return

	} else if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) CreateLivechatInteraction(c *gin.Context) {
	var clcir presentation.CreateLiveChatInteractionRequest
	errorMessage := make(map[string]string)

	err := c.BindJSON(&clcir)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Create Live Chat Interaction] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, err := ih.interactionService.CreateLiveChatInteraction(&clcir)
	if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (ih *InteractionHandler) GetGeotagInformation(c *gin.Context) {
	var geoTag presentation.GetGeotagInformation
	errorMessage := make(map[string]string)

	err := c.BindJSON(&geoTag)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Geotag Information] Invalid Body : %v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, err := ih.interactionService.GetGeotagInformation(geoTag.InterractionId, geoTag)
	if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}
