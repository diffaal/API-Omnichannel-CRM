package handler

import (
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/presentation"
	"Omnichannel-CRM/package/response"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChannelAccountHandler struct {
	channelAccountService service.IChannelAccountService
}

func NewChannelAccountHandler(channelAccountService service.IChannelAccountService) *ChannelAccountHandler {
	channelAccountHandler := ChannelAccountHandler{
		channelAccountService: channelAccountService,
	}
	return &channelAccountHandler
}

func (cah *ChannelAccountHandler) CreateChannelAccount(c *gin.Context) {
	var cam presentation.ChannelAccountModel
	errorMessage := make(map[string]string)

	err := c.BindJSON(&cam)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Create Channel Account] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, err := cah.channelAccountService.CreateChannelAccount(&cam)
	if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}

func (cah *ChannelAccountHandler) GetChannelAccountList(c *gin.Context) {
	errorMessage := make(map[string]string)

	result, err := cah.channelAccountService.GetChannelAccountList()
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

func (cah *ChannelAccountHandler) GetChannelAccountById(c *gin.Context) {
	errorMessage := make(map[string]string)

	channelAccountIdQuery := c.Query("channel_account_id")
	if channelAccountIdQuery == "" {
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	channelAccountId, err := strconv.ParseUint(channelAccountIdQuery, 10, 64)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Channel Account] Invalid Value of Query channel_account_id"))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}
	caIdUint := uint(channelAccountId)

	result, err := cah.channelAccountService.GetChannelAccountById(caIdUint)
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

func (cah *ChannelAccountHandler) UpdateChannelAccount(c *gin.Context) {
	var ucam presentation.UpdateChannelAccountModel
	errorMessage := make(map[string]string)

	err := c.BindJSON(&ucam)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Update Channel Account] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, err := cah.channelAccountService.UpdateChannelAccount(&ucam)
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

func (cah *ChannelAccountHandler) DeleteChannelAccount(c *gin.Context) {
	var dcam presentation.DeleteChannelAccountModel
	errorMessage := make(map[string]string)

	err := c.BindJSON(&dcam)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Delete Channel Account] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, err := cah.channelAccountService.DeleteChannelAccountById(&dcam)
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
