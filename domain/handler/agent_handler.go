package handler

import (
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/presentation"
	"Omnichannel-CRM/package/response"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentService service.IAgentService
}

func NewAgentHandler(agentService service.IAgentService) *AgentHandler {
	agentHandler := AgentHandler{
		agentService: agentService,
	}
	return &agentHandler
}

func (ah *AgentHandler) GetDashboardAgentList(c *gin.Context) {
	errorMessage := make(map[string]string)

	filters, err := presentation.ParseGetListAgentFilters(c)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Agent List] Invalid Query Params: %+v", err))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	result, err := ah.agentService.GetAgentList(filters)
	if errors.Is(err, enum.ERROR_DATA_NOT_FOUND) {
		errorMessage["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		errorMessage["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		response.ResponseNotFound(c, nil, errorMessage)
		return

	} else if err != nil {
		errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		logger.Info(fmt.Sprintf("[FAILED][Get Dashboard Agent List] Internal Server Error: %+v", err))
		response.ResponseInternalServerError(c, nil, errorMessage)
		return
	}

	response.ResponseWithData(c, result, errorMessage)
}
