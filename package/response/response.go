package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseBase struct {
	Data         interface{} `json:"data"`
	IsError      bool        `json:"isError"`
	ErrorStatus  string      `json:"errorStatus"`
	ErrorMessage string      `json:"errorMessage"`
}

func ResponseWithData(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusOK, response)
}

func ResponseInvalidRequest(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusUnprocessableEntity, response)
}

func ResponseNotFound(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusNotFound, response)
}

func ResponseBadRequest(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusBadRequest, response)
}

func ResponseUnauthorized(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusUnauthorized, response)
}

func ResponseForbidden(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusForbidden, response)
}

func ResponseInternalServerError(c *gin.Context, data interface{}, errorMessage map[string]string) {
	isError := false

	if errorMessage["errorStatus"] != "" {
		isError = true
	}

	response := ResponseBase{
		IsError:      isError,
		Data:         data,
		ErrorStatus:  errorMessage["errorStatus"],
		ErrorMessage: errorMessage["errorMessage"],
	}

	c.IndentedJSON(http.StatusInternalServerError, response)
}
