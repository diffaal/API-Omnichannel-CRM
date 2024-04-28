package api

import (
	"Omnichannel-CRM/api/middleware"
	"Omnichannel-CRM/domain/handler"
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/logger"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/gmail/v1"
	"gorm.io/gorm"
)

func SetupRouter(dbCRM *gorm.DB, dbOmnichannel *gorm.DB) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(gin.LoggerWithWriter(logger.Logger.Writer()))

	router.Static("/files", "/var/www")

	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	interactionRepo := repository.NewInteractionRepository(dbOmnichannel)
	messageRepo := repository.NewMessageRepository(dbOmnichannel)
	userRepo := repository.NewUserRepository(dbCRM)
	reporterRepo := repository.NewReporterRepository(dbOmnichannel)

	gmailService := service.NewGmailService()
	emailRepo := repository.NewEmailRepository(dbOmnichannel, gmailService)
	threadRepo := repository.NewThreadRepository(dbOmnichannel)

	emailService := service.NewEmailService(interactionRepo, messageRepo, reporterRepo, emailRepo, *threadRepo)

	interactionService := service.NewInteractionService(interactionRepo, messageRepo, userRepo, reporterRepo, emailService, threadRepo)
	interactionHandler := handler.NewInteractionHandler(interactionService)

	interactionApi := router.Group("interaction/")
	{
		interactionApi.GET("/list", middleware.AuthMiddleware(), interactionHandler.GetDashboardInteractionList)
		interactionApi.PUT("/claim", middleware.AuthMiddleware(), interactionHandler.ClaimInteractionByAgent)
		interactionApi.POST("/messenger/send", middleware.AuthMiddleware(), interactionHandler.MessengerSendMessage)
		interactionApi.POST("/live-chat/send", interactionHandler.LiveChatSendMessage)
		interactionApi.GET("/messages", interactionHandler.GetInteractionMessages)
		interactionApi.GET("/my", middleware.AuthMiddleware(), interactionHandler.GetAgentInteractions)
		interactionApi.PUT("/close", middleware.AuthMiddleware(), interactionHandler.CloseInteractionByAgent)
		interactionApi.GET("/closed-data", middleware.AuthMiddleware(), interactionHandler.GetClosedInteractionsData)
		interactionApi.POST("/live-chat/create", interactionHandler.CreateLivechatInteraction)
	}

	router.POST("/geotag", interactionHandler.GetGeotagInformation)

	reporterService := service.NewReporterService(reporterRepo)
	reporterHandler := handler.NewReporterHandler(reporterService)

	reporterApi := router.Group("reporter/")
	{
		reporterApi.GET("/get", reporterHandler.GetReporterByReporterId)
		reporterApi.PUT("/update", reporterHandler.UpdateReporter)
	}

	agentService := service.NewAgentService(userRepo, interactionRepo)
	agentHandler := handler.NewAgentHandler(agentService)

	agentApi := router.Group("/agent")
	{
		agentApi.GET("/list", agentHandler.GetDashboardAgentList)
	}

	channelAccountRepo := repository.NewChannelAccountRepository(dbCRM)
	channelAccountService := service.NewChannelAccountService(channelAccountRepo)
	channelAccountHandler := handler.NewChannelAccountHandler(channelAccountService)
	channelAccountApi := router.Group("/channel-account")
	{
		channelAccountApi.POST("/create", channelAccountHandler.CreateChannelAccount)
		channelAccountApi.GET("/list", channelAccountHandler.GetChannelAccountList)
		channelAccountApi.GET("/get", channelAccountHandler.GetChannelAccountById)
		channelAccountApi.PUT("/update", channelAccountHandler.UpdateChannelAccount)
		channelAccountApi.DELETE("/delete", channelAccountHandler.DeleteChannelAccount)
	}

	return router
}

func SetupWebhookRouter(dbOmnichannel *gorm.DB) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(gin.LoggerWithWriter(logger.Logger.Writer()))

	router.Static("/files", "/var/www")

	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	interactionRepo := repository.NewInteractionRepository(dbOmnichannel)
	messageRepo := repository.NewMessageRepository(dbOmnichannel)
	reporterRepo := repository.NewReporterRepository(dbOmnichannel)
	metaWebhookService := service.NewMetaWebhookService(interactionRepo, messageRepo, reporterRepo)
	metaWebhookHandler := handler.NewMetaWebhookHandler(metaWebhookService)

	gmailService := service.NewGmailService()
	emailRepo := repository.NewEmailRepository(dbOmnichannel, gmailService)
	threadRepo := repository.NewThreadRepository(dbOmnichannel)
	emailService := service.NewEmailService(interactionRepo, messageRepo, reporterRepo, emailRepo, *threadRepo)

	watchRes, err := gmailService.Users.Watch("me", &gmail.WatchRequest{
		LabelIds:  []string{"INBOX", "UNREAD"},
		TopicName: "projects/email-api-406816/topics/gmail-webhook",
	}).Do()
	if err != nil {
		fmt.Println("Error when creating watch request to google api")
	}

	var gmailHandler handler.GmailHandler
	if watchRes != nil {
		gmailHandler = handler.NewGmailHandler(emailService, int(watchRes.HistoryId))
	} else {
		gmailHandler = handler.NewGmailHandler(emailService, 0)
	}

	metaWebhookApi := router.Group("webhook-meta/")
	{
		metaWebhookApi.POST("/facebook", metaWebhookHandler.FacebookWebhookInteractionHandler)
		metaWebhookApi.GET("/facebook", metaWebhookHandler.ValidateVerificationRequest)
		metaWebhookApi.POST("/instagram", metaWebhookHandler.InstagramWebhookInteractionHandler)
		metaWebhookApi.GET("/instagram", metaWebhookHandler.ValidateVerificationRequest)
	}

	whatsappWebhookApi := router.Group("/webhooks")
	{
		whatsappWebhookApi.POST("", metaWebhookHandler.WhatsappMessageInteractionHandler)
		whatsappWebhookApi.GET("", metaWebhookHandler.ValidateVerificationRequest)
	}

	gmailWebhook := router.Group("/webhook-gmail")
	{
		gmailWebhook.POST("", gmailHandler.Webhook)
	}

	return router
}
