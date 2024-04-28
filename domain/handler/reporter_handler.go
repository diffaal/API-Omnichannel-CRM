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

type ReporterHandler struct {
	reporterService service.IReporterService
}

func NewReporterHandler(reporterService service.IReporterService) *ReporterHandler {
	reporterHandler := ReporterHandler{
		reporterService: reporterService,
	}
	return &reporterHandler
}

func (rh *ReporterHandler) GetReporterByReporterId(c *gin.Context) {
	errorMessage := make(map[string]string)

	reporterIdQuery := c.Query("reporter_id")
	if reporterIdQuery == "" {
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}

	reporterId, err := strconv.ParseUint(reporterIdQuery, 10, 64)
	if err != nil {
		errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
		errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Get Reporter] Invalid Value of Query reporter_id"))
		response.ResponseInvalidRequest(c, nil, errorMessage)
		return
	}
	reporterIdUint := uint(reporterId)

	result, err := rh.reporterService.GetReporterByReporterId(reporterIdUint)
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

func (rh *ReporterHandler) UpdateReporter(c *gin.Context) {
	var urr presentation.UpdateReporterRequest
	errorMessage := make(map[string]string)

	err := c.BindJSON(&urr)
	if err != nil {
		errorMessage["errorMessage"] = enum.FAILED_BIND_JSON_MESSAGE
		errorMessage["errorStatus"] = enum.FAILED_BIND_JSON_STATUS
		logger.Info(fmt.Sprintf("[FAILED][Update Reporter] Bind JSON Body: %+v", err))
		response.ResponseBadRequest(c, nil, errorMessage)
		return
	}

	result, err := rh.reporterService.UpdateReporter(&urr)
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
