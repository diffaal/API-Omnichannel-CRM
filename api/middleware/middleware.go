package middleware

import (
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/jwt"
	"Omnichannel-CRM/package/response"

	"Omnichannel-CRM/package/logger"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorMessage := make(map[string]string)
		var atn jwt.AccessTokenNodes
		err := jwt.TokenValid(c.Request)
		if err != nil {
			errorMessage := map[string]string{
				"errorStatus":  enum.UNAUTHORIZED_STATUS,
				"errorMessage": enum.UNAUTHORIZED_MESSAGE,
			}
			response.ResponseUnauthorized(c, "", errorMessage)
			c.Abort()
			return
		}

		atn.AccessToken = jwt.ExtractToken(c.Request)

		// Extract Jwt String
		responseUserData := jwt.GetDataFromToken(&atn)

		parseValueUserId := responseUserData["user_id"].(string)
		if parseValueUserId == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The user_id filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueUsername := responseUserData["username"].(string)
		if parseValueUsername == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The username filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueRole := int(responseUserData["role"].(float64))
		if parseValueRole == 0 {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The role filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueChannelAccount := responseUserData["channel_account"]

		c.Set("user_id", parseValueUserId)
		c.Set("username", parseValueUsername)
		c.Set("role", parseValueRole)
		c.Set("channel_account", parseValueChannelAccount)

		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorMessage := make(map[string]string)
		var atn jwt.AccessTokenNodes
		err := jwt.TokenValid(c.Request)
		if err != nil {
			errorMessage := map[string]string{
				"errorStatus":  enum.UNAUTHORIZED_STATUS,
				"errorMessage": enum.UNAUTHORIZED_MESSAGE,
			}
			response.ResponseUnauthorized(c, "", errorMessage)
			c.Abort()
			return
		}
		atn.AccessToken = jwt.ExtractToken(c.Request)

		// Extract Jwt String
		responseUserData := jwt.GetDataFromToken(&atn)

		parseValueUserId := responseUserData["user_id"].(string)
		if parseValueUserId == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The user_id filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueUsername := responseUserData["username"].(string)
		if parseValueUsername == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The username filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueRole := int(responseUserData["role"].(float64))
		if parseValueRole == 0 {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The role filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		if parseValueRole != enum.ROLE_ADMIN_PUSAT && parseValueRole != enum.ROLE_ADMIN_PROVINSI && parseValueRole != enum.ROLE_ADMIN_KOTA {
			errorMessage["errorMessage"] = enum.FORBIDDEN_MESSAGE
			errorMessage["errorStatus"] = enum.FORBIDDEN_STATUS
			logger.Info("[FAILED][AdminAuthMiddleware] The role is forbidden to proceed")
			response.ResponseForbidden(c, nil, errorMessage)
			c.Abort()
			return
		}

		c.Set("user_id", parseValueUserId)
		c.Set("username", parseValueUsername)
		c.Set("role", parseValueRole)

		c.Next()
	}
}

func DispatcherAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorMessage := make(map[string]string)
		var atn jwt.AccessTokenNodes
		err := jwt.TokenValid(c.Request)
		if err != nil {
			errorMessage := map[string]string{
				"errorStatus":  enum.UNAUTHORIZED_STATUS,
				"errorMessage": enum.UNAUTHORIZED_MESSAGE,
			}
			response.ResponseUnauthorized(c, "", errorMessage)
			c.Abort()
			return
		}

		atn.AccessToken = jwt.ExtractToken(c.Request)

		// Extract Jwt String
		responseUserData := jwt.GetDataFromToken(&atn)

		parseValueUserId := responseUserData["user_id"].(string)
		if parseValueUserId == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The user_id filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueUsername := responseUserData["username"].(string)
		if parseValueUsername == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The username filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueRole := int(responseUserData["role"].(float64))
		if parseValueRole == 0 {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The role filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		if parseValueRole != enum.ROLE_DISPATCHER_PROVINSI && parseValueRole != enum.ROLE_ADMIN_KOTA {
			errorMessage["errorMessage"] = enum.FORBIDDEN_MESSAGE
			errorMessage["errorStatus"] = enum.FORBIDDEN_STATUS
			logger.Info("[FAILED][AdminAuthMiddleware] The role is forbidden to proceedd")
			response.ResponseForbidden(c, nil, errorMessage)
			c.Abort()
			return
		}

		c.Set("user_id", parseValueUserId)
		c.Set("username", parseValueUsername)
		c.Set("role", parseValueRole)

		c.Next()
	}
}

func AgentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorMessage := make(map[string]string)
		var atn jwt.AccessTokenNodes
		err := jwt.TokenValid(c.Request)
		if err != nil {
			errorMessage := map[string]string{
				"errorStatus":  enum.UNAUTHORIZED_STATUS,
				"errorMessage": enum.UNAUTHORIZED_MESSAGE,
			}
			response.ResponseUnauthorized(c, "", errorMessage)
			c.Abort()
			return
		}

		atn.AccessToken = jwt.ExtractToken(c.Request)

		// Extract Jwt String
		responseUserData := jwt.GetDataFromToken(&atn)

		parseValueUserId := responseUserData["user_id"].(string)
		if parseValueUserId == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The user_id filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueUsername := responseUserData["username"].(string)
		if parseValueUsername == "" {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The username filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		parseValueRole := int(responseUserData["role"].(float64))
		if parseValueRole == 0 {
			errorMessage["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			errorMessage["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			logger.Info("[FAILED][Get Profile] The role filled not found")
			response.ResponseWithData(c, nil, errorMessage)
			c.Abort()
			return
		}

		if parseValueRole != enum.ROLE_DISPATCHER_KOTA && parseValueRole != enum.ROLE_RESPONDER_KOTA {
			errorMessage["errorMessage"] = enum.FORBIDDEN_MESSAGE
			errorMessage["errorStatus"] = enum.FORBIDDEN_STATUS
			logger.Info("[FAILED][AdminAuthMiddleware] The role is forbidden to proceedd")
			response.ResponseForbidden(c, nil, errorMessage)
			c.Abort()
			return
		}

		c.Set("user_id", parseValueUserId)
		c.Set("username", parseValueUsername)
		c.Set("role", parseValueRole)

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Method", "POST, GET, DELETE, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
