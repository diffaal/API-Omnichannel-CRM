package service

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/package/config"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/presentation"
	"Omnichannel-CRM/package/request"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	//"sort"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type InteractionService struct {
	interactionRepo repository.IinteractionRepository
	messageRepo     repository.IMessageRepository
	userRepo        repository.IUserRepository
	reporterRepo    repository.IReporterRepository
	emailService    IEmailService
	threadRepo      repository.IThreadRepository
}

type IInteractionService interface {
	UpdateInteractionStatusByAgent(*presentation.ClaimInteractionRequest, string, string) (*entity.Interaction, error)
	GetInteractionList(map[string]interface{}, *entity.ChannelAccount) (map[string]interface{}, error)
	GetInteractionMessages(uint) (map[string]interface{}, error)
	GetAgentInteractions(string, map[string]interface{}) (map[string]interface{}, error)

	MessengerSendMessagetoMeta(*presentation.MetaSendMessageRequest, *entity.ChannelAccount) (map[string]interface{}, *entity.Message, error)
	WhatsappSendMessagetoMeta(*presentation.MetaSendMessageRequest, *entity.ChannelAccount) (map[string]interface{}, *entity.Message, error)
	LiveChatSendMessage(*presentation.MetaSendMessageRequest) (map[string]interface{}, *entity.Message, error)

	GetClosedInteractionsData() (map[string]interface{}, error)
	SendClosedInteractionData(*entity.Interaction) error
	SendMessageToEmail(req presentation.MessengerSendEmailRequest) (map[string]interface{}, *entity.Message, error)
	CreateLiveChatInteraction(*presentation.CreateLiveChatInteractionRequest) (map[string]interface{}, error)

	GetGeotagInformation(uint, presentation.GetGeotagInformation) (entity.GeotagInformation, error)

	WebsocketSendService(messages entity.Message) error
}

func NewInteractionService(interactionRepo repository.IinteractionRepository, messageRepo repository.IMessageRepository, userRepo repository.IUserRepository, reporterRepo repository.IReporterRepository, emailService IEmailService, threadRepo repository.IThreadRepository) *InteractionService {
	interactionService := InteractionService{
		interactionRepo: interactionRepo,
		messageRepo:     messageRepo,
		userRepo:        userRepo,
		reporterRepo:    reporterRepo,
		emailService:    emailService,
		threadRepo:      threadRepo,
	}
	return &interactionService
}

func (is *InteractionService) GetGeotagInformation(interactionId uint, address presentation.GetGeotagInformation) (entity.GeotagInformation, error) {
	geoTagInfo := []entity.GeotagInformation{}
	lat, lon := address.Lat, address.Lon
	if lat == "" || lon == "" {
		posturl := "https://geocode.maps.co/search?q=%s&api_key=65a748960f046188097979vzq9fd150"
		url := fmt.Sprintf(posturl, address.Address)

		r, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return entity.GeotagInformation{}, err
		}

		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			return entity.GeotagInformation{}, err
		}

		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&geoTagInfo)
		if err != nil {
			return entity.GeotagInformation{}, err
		}

		if len(geoTagInfo) == 0 {
			return entity.GeotagInformation{}, enum.ERROR_DATA_NOT_FOUND
		}

		lat = geoTagInfo[0].Lat
		lon = geoTagInfo[0].Lon
	}

	interaction := entity.Interaction{
		Latitude:  lat,
		Longitude: lon,
	}

	_, err := is.interactionRepo.UpdateInteraction(interactionId, &interaction)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.GeotagInformation{}, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return entity.GeotagInformation{}, err
	}

	info := entity.GeotagInformation{Lat: lat, Lon: lon}

	return info, nil

}

var connections = make(map[string]*WebSocketClient)

func (mws *InteractionService) WebsocketSendService(message entity.Message) error {
	config.GetConfig()

	host := viper.GetString("Websocket.Host")
	channel := "ws"

	query := fmt.Sprintf("user_id=%s&room_id=%d", message.SenderId, message.InteractionId)
	if message.SentBy == enum.AGENT {
		query += "&is_agent=true"
	}
	client, err := NewWebSocketClient(host, channel, fmt.Sprint(message.InteractionId), query)
	if err != nil {
		return err
	}

	send := MessageToSend{
		Action: "send-message",
		Message: presentation.Message{
			ID:               message.ID,
			CreatedAt:        message.CreatedAt,
			UpdatedAt:        message.UpdatedAt,
			InteractionId:    message.InteractionId,
			SenderId:         message.SenderId,
			RecipientId:      message.RecipientId,
			MetaMessageId:    message.MetaMessageId,
			Message:          message.Message,
			MessageTimestamp: message.MessageTimestamp,
			SentBy:           message.SentBy,
			IsRead:           message.IsRead,
			IsDeleted:        message.IsDeleted,
		},
	}
	err = client.Write(send)
	if err != nil {
		return err
	}

	return nil
}

func (is *InteractionService) GetInteractionList(filters map[string]interface{}, channelAccount *entity.ChannelAccount) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var dil []presentation.DashboardInteractionList
	var agentIds []string
	var uniqueAgentIds []string
	var agentList []entity.User
	var count int64
	checkDuplicateAgentIds := make(map[string]bool)

	interactionList, count, err := is.interactionRepo.GetInteractionList(filters, channelAccount)
	if (interactionList == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		result["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		result["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		return result, err

	} else if err != nil {
		result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		return result, err
	}

	for _, v := range interactionList {
		agentIds = append(agentIds, v.AgentId)
	}

	for _, aid := range agentIds {
		if _, v := checkDuplicateAgentIds[aid]; !v {
			checkDuplicateAgentIds[aid] = v
			uniqueAgentIds = append(uniqueAgentIds, aid)
		}
	}

	agentList, err = is.userRepo.GetUserListByIds(uniqueAgentIds)
	if err != nil {
		result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		return result, err
	}

	for _, v := range interactionList {
		dild := presentation.DashboardInteractionList{
			InteractionId:   v.ID,
			PlatformId:      v.PlatformId,
			ReporterId:      v.ReporterId,
			ConversationId:  v.ConversationId,
			MentionMediaId:  v.MentionMediaId,
			AgentId:         v.AgentId,
			Status:          v.Status,
			Platform:        v.Platform,
			InteractionType: v.InteractionType,
			Duration:        v.Duration,
		}
		for _, w := range agentList {
			if w.ID == v.AgentId {
				dild.AgentName = fmt.Sprintf("%s %s", w.FirstName, w.LastName)
			}
		}
		dil = append(dil, dild)
	}

	result["interaction_list"] = dil
	result["page"] = filters["page"]
	result["pageSize"] = filters["pageSize"]
	result["total"] = count

	return result, nil
}

func (is *InteractionService) UpdateInteractionStatusByAgent(cir *presentation.ClaimInteractionRequest, userId string, status string) (*entity.Interaction, error) {
	interactionId := cir.InteractionId
	newInteraction := entity.Interaction{
		AgentId: userId,
		Status:  status,
	}

	interaction, err := is.interactionRepo.UpdateInteraction(interactionId, &newInteraction)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	return interaction, nil
}

func (is *InteractionService) LiveChatSendMessage(msmr *presentation.MetaSendMessageRequest) (map[string]interface{}, *entity.Message, error) {
	result := make(map[string]interface{})

	interaction, err := is.interactionRepo.GetInteractionById(msmr.InteractionId)
	if (interaction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		result["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		result["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		return result, nil, err

	} else if err != nil {
		result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		return result, nil, err
	}

	var senderId string
	var recipientId string

	if msmr.SentBy == enum.AGENT {
		senderId = msmr.PlatformId
		recipientId = fmt.Sprint(msmr.ReporterId)
	} else {
		senderId = fmt.Sprint(msmr.ReporterId)
		recipientId = msmr.PlatformId
	}

	sendedMessage := entity.Message{
		InteractionId:    msmr.InteractionId,
		SenderId:         senderId,
		RecipientId:      recipientId,
		Message:          msmr.Message,
		MessageTimestamp: time.Now(),
		SentBy:           msmr.SentBy,
		IsRead:           false,
	}

	message, err := is.messageRepo.CreateMessage(&sendedMessage)
	if err != nil {
		return nil, nil, err
	}

	result["message_id"] = message.ID

	return result, message, nil
}

func (is *InteractionService) MessengerSendMessagetoMeta(msmr *presentation.MetaSendMessageRequest, channelAccount *entity.ChannelAccount) (map[string]interface{}, *entity.Message, error) {
	result := make(map[string]interface{})
	var access_token string

	if channelAccount.ID == 0 {
		return nil, nil, enum.USER_DO_NOT_HAVE_CHANNEL_ACCOUNT
	}

	if channelAccount.FaceboookPageId == "" {
		return nil, nil, enum.PLATFORM_ID_NOT_SET
	}

	reporter, err := is.reporterRepo.GetReporterByReporterId(msmr.ReporterId)
	if err != nil {
		return nil, nil, err
	}

	platform := msmr.Platform

	if platform == enum.FACEBOOK {
		if msmr.PlatformId != channelAccount.FaceboookPageId {
			return nil, nil, enum.CHANNEL_ACCOUNT_NOT_MATCH
		}
		if channelAccount.FacebookAccessToken == "" {
			return nil, nil, enum.PLATFORM_ACCESS_TOKEN_NOT_SET
		}
		access_token = channelAccount.FacebookAccessToken
	} else {
		if channelAccount.InstagramId == "" {
			return nil, nil, enum.PLATFORM_ID_NOT_SET
		}

		if msmr.PlatformId != channelAccount.InstagramId {
			return nil, nil, enum.CHANNEL_ACCOUNT_NOT_MATCH
		}
		if channelAccount.InstagramAccessToken == "" {
			return nil, nil, enum.PLATFORM_ACCESS_TOKEN_NOT_SET
		}
		access_token = channelAccount.InstagramAccessToken
	}

	path := fmt.Sprintf("/%s/%s/messages", viper.GetString("Meta.API_VERSION"), channelAccount.FaceboookPageId)
	params := url.Values{}
	params.Add("access_token", access_token)

	reqUrl := url.URL{
		Scheme:   "https",
		Host:     "graph.facebook.com",
		Path:     path,
		RawQuery: params.Encode(),
	}

	body := presentation.MessengerSendMessageMetaRequest{
		Recipient:     presentation.IdField{Id: reporter.MetaReporterId},
		MessagingType: "MESSAGE_TAG",
		Tag:           "HUMAN_AGENT",
		Message:       presentation.MessageSendMetaField{MessageText: msmr.Message},
	}

	response, err := request.PostRequest(reqUrl, body, "")
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		messageData := &presentation.MessengerSendMessageMetaResponse{}
		err = json.NewDecoder(response.Body).Decode(messageData)
		if err != nil {
			return nil, nil, err
		}

		sendedMessage := entity.Message{
			InteractionId:    msmr.InteractionId,
			SenderId:         msmr.PlatformId,
			RecipientId:      messageData.ReporterId,
			MetaMessageId:    messageData.MessageId,
			Message:          msmr.Message,
			MessageTimestamp: time.Now(),
			SentBy:           enum.AGENT,
			IsRead:           false,
		}

		message, err := is.messageRepo.CreateMessage(&sendedMessage)
		if err != nil {
			return nil, nil, err
		}

		result["message_id"] = message.MetaMessageId

		return result, message, nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	bodyString := string(bodyBytes)
	logger.Info(response.Status)
	logger.Info(bodyString)
	return nil, nil, errors.New("meta response not 200")
}

func (is *InteractionService) WhatsappSendMessagetoMeta(msmr *presentation.MetaSendMessageRequest, channelAccount *entity.ChannelAccount) (map[string]interface{}, *entity.Message, error) {
	result := make(map[string]interface{})

	if channelAccount.ID == 0 {
		return nil, nil, enum.USER_DO_NOT_HAVE_CHANNEL_ACCOUNT
	}

	if channelAccount.WhatsappNumId == "" {
		return nil, nil, enum.PLATFORM_ID_NOT_SET
	}

	if channelAccount.WhatsappBusinessId == "" {
		return nil, nil, enum.PLATFORM_ID_NOT_SET
	}

	if msmr.PlatformId != channelAccount.WhatsappBusinessId {
		return nil, nil, enum.CHANNEL_ACCOUNT_NOT_MATCH
	}

	if channelAccount.WhatsappAccessToken == "" {
		return nil, nil, enum.PLATFORM_ACCESS_TOKEN_NOT_SET
	}
	access_token := channelAccount.WhatsappAccessToken

	reporter, err := is.reporterRepo.GetReporterByReporterId(msmr.ReporterId)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/%s/%s/messages", viper.GetString("Meta.WA_API_VERSION"), channelAccount.WhatsappNumId)

	reqUrl := url.URL{
		Scheme: "https",
		Host:   "graph.facebook.com",
		Path:   path,
	}

	body := presentation.WhatsappSendMessageMetaRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		Recipient:        reporter.MetaReporterId,
		MessageType:      "text",
		Text: presentation.WhatsappMessageMetaField{
			PreviewUrl: false,
			Body:       msmr.Message,
		},
	}

	response, err := request.PostRequest(reqUrl, body, access_token)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		messageData := &presentation.WhatsappSendMessageMetaResponse{}
		err = json.NewDecoder(response.Body).Decode(messageData)
		if err != nil {
			return nil, nil, err
		}

		sendedMessage := entity.Message{
			InteractionId:    msmr.InteractionId,
			SenderId:         msmr.PlatformId,
			RecipientId:      messageData.Contacts[0].WaId,
			MetaMessageId:    messageData.Messages[0].Id,
			Message:          msmr.Message,
			MessageTimestamp: time.Now(),
			SentBy:           enum.AGENT,
			IsRead:           false,
		}

		message, err := is.messageRepo.CreateMessage(&sendedMessage)
		if err != nil {
			return nil, nil, err
		}

		result["message_id"] = message.MetaMessageId

		return result, message, nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	bodyString := string(bodyBytes)
	logger.Info(response.Status)
	logger.Info(bodyString)
	return nil, nil, errors.New("meta response not 200")
}

func (is *InteractionService) GetInteractionMessages(interactionId uint) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	interaction, err := is.interactionRepo.GetInteractionById(interactionId)
	if (interaction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		result["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		result["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		return result, err

	} else if err != nil {
		result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		return result, err
	}

	messages, err := is.messageRepo.GetMessagesofInteraction(interactionId)
	if (interaction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		result["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		result["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		return result, err

	} else if err != nil {
		result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		return result, err
	}

	if interaction.Platform == enum.EMAIL {
		thread, err := is.threadRepo.GetThreadByID(interaction.ConversationId)
		if err != nil {
			result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
			result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
			return result, err
		}

		result["threadInfo"] = thread
	}

	result["interaction"] = interaction
	result["messages"] = messages

	return result, nil
}

func (is *InteractionService) GetAgentInteractions(userId string, filters map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	interactions, err := is.interactionRepo.GetInteractionsByAgentId(userId, filters)
	if (interactions == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		result["errorStatus"] = enum.DATA_NOT_FOUND_STATUS
		result["errorMessage"] = enum.DATA_NOT_FOUND_MESSAGE
		return result, err

	} else if err != nil {
		result["errorStatus"] = enum.SYSTEM_BUSY_STATUS
		result["errorMessage"] = enum.SYSTEM_BUSY_MESSAGE
		return result, err
	}

	result["interaction_list"] = interactions

	return result, nil
}

func (is *InteractionService) GetClosedInteractionsData() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	filter := make(map[string]interface{})
	filter["status"] = []string{enum.CLOSED}
	dataList := []presentation.InteractionReporterData{}
	channelAccount := entity.ChannelAccount{}

	interactionList, _, err := is.interactionRepo.GetInteractionList(filter, &channelAccount)
	if (interactionList == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	for _, v := range interactionList {
		data := presentation.InteractionReporterData{Interaction: v}
		reporter, err := is.reporterRepo.GetReporterByReporterId(v.ReporterId)
		if (interactionList == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
			dataList = append(dataList, data)
			continue

		} else if err != nil {
			return nil, err
		}

		data.Reporter = *reporter
		dataList = append(dataList, data)
	}

	result["data_list"] = dataList

	return result, nil
}

func (is *InteractionService) SendClosedInteractionData(interaction *entity.Interaction) error {
	var body presentation.SendInteractionDataCRMRequest

	if interaction.InteractionType != enum.MENTION {
		reporter, err := is.reporterRepo.GetReporterByReporterId(interaction.ReporterId)
		if err != nil {
			return err
		}
		body.ReporterName = reporter.Name
		body.ReporterGender = reporter.Gender
		body.ReporterPhoneNumber = reporter.PhoneNumber
		body.ReporterEmail = reporter.Email
		body.ReporterAddress = reporter.Address
	}

	access_token := viper.GetString("SERVER_TOKEN")

	reqUrl := url.URL{
		Scheme: "http",
		Host:   viper.GetString("Services.CRM_URL"),
		Path:   "/interaction",
	}

	body.InteractionId = interaction.ID
	body.CreatedAt = interaction.CreatedAt
	body.UpdatedAt = interaction.UpdatedAt
	body.DeletedAt = interaction.DeletedAt.Time
	body.PlatformId = interaction.PlatformId
	body.ConversationId = interaction.ConversationId
	body.ReporterId = interaction.ReporterId
	body.MentionMediaId = interaction.MentionMediaId
	body.MentionMediaUrl = interaction.MentionMediaUrl
	body.AgentId = interaction.AgentId
	body.Status = interaction.Status
	body.Platform = interaction.Platform
	body.InteractionType = interaction.InteractionType
	body.Duration = interaction.Duration
	body.Latitude = interaction.Latitude
	body.Longitude = interaction.Longitude

	response, err := request.PostRequest(reqUrl, body, access_token)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	logger.Info(response.Status)
	logger.Info(bodyString)

	respBody := make(map[string]interface{})

	err = json.Unmarshal(bodyBytes, &respBody)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%v", respBody))

	if respBody["isError"] == true {
		logger.Info(fmt.Sprintf("%v", respBody["errorStatus"]))
		logger.Info(fmt.Sprintf("%v", respBody["errorMessage"]))
		return enum.CRM_RESPONSE_ERROR
	}

	return nil
}

func (is *InteractionService) SendMessageToEmail(req presentation.MessengerSendEmailRequest) (map[string]interface{}, *entity.Message, error) {
	messageId, message, err := is.emailService.SendEmail(req.InteractionId, req.Message)
	if err != nil {
		return nil, nil, fmt.Errorf("[InteractionService][SendMessageToEmail] Error when calling SendEmail, trace: %+v", err)
	}

	res := map[string]interface{}{
		"message_id": messageId,
	}

	return res, message, nil
}

func (is *InteractionService) CreateLiveChatInteraction(clcir *presentation.CreateLiveChatInteractionRequest) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	newReporter := entity.Reporter{
		Name:        clcir.Name,
		Email:       clcir.Email,
		PhoneNumber: clcir.PhoneNumber,
	}

	reporter, err := is.reporterRepo.CreateReporter(&newReporter)
	if err != nil {
		return nil, err
	}

	newInteraction := entity.Interaction{
		ReporterId:      reporter.ID,
		Status:          enum.UNCLAIMED,
		Platform:        enum.LIVE_CHAT,
		InteractionType: enum.PESAN,
	}

	interaction, err := is.interactionRepo.CreateInteraction(&newInteraction)
	if err != nil {
		return nil, err
	}

	result["reporter"] = reporter
	result["interaction"] = interaction

	return result, nil
}
