package service

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/package/config"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"gorm.io/gorm"
)

type EmailService struct {
	interactionRepo repository.IinteractionRepository
	messageRepo     repository.IMessageRepository
	reporterRepo    repository.IReporterRepository
	emailRepo       repository.IEmailRepository
	threadRepo      repository.ThreadRepository
}

type IEmailService interface {
	ProcessWebhook(rawMessage string, prevHistoryId int) (historyId uint64, err error)
	SendEmail(interactionId uint, message string) (messageId string, res *entity.Message, err error)
}

func init() {
	config.GetConfig()
}

func NewEmailService(interactionRepo repository.IinteractionRepository, messageRepo repository.IMessageRepository, reporterRepo repository.IReporterRepository, emailRepo repository.IEmailRepository, threadRepo repository.ThreadRepository) *EmailService {
	emailService := EmailService{
		interactionRepo: interactionRepo,
		messageRepo:     messageRepo,
		// userRepo:        userRepo,
		reporterRepo: reporterRepo,
		emailRepo:    emailRepo,
		threadRepo:   threadRepo,
	}
	return &emailService
}

type ResponseData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func RefreshToken() (accessToken string, err error) {
	refreshToken := viper.GetString("OAuth.Refresh_Token")
	clientSecret := viper.GetString("OAuth.Client_Secret")
	clientId := viper.GetString("OAuth.Client_Id")

	reqUrl := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("refresh_token", refreshToken)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Define a struct for the JSON response
	var responseData ResponseData

	// Unmarshal JSON into struct
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", err
	}

	return responseData.AccessToken, nil
}

func NewGmailService() *gmail.Service {
	accessToken := viper.GetString("OAuth.Access_Token")

	accessToken, err := RefreshToken()
	if err != nil {
		accessToken = viper.GetString("OAuth.Access_Token")
	}
	// Create an OAuth2 token source using the access token
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	// Create an OAuth2 HTTP client from the token source
	oauth2Client := oauth2.NewClient(context.Background(), tokenSource)

	// Create a Gmail service using the OAuth2 client
	gmailService, err := gmail.New(oauth2Client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	return gmailService
}

func FindHeaders(headers []*gmail.MessagePartHeader, key string) (value string) {
	for _, elem := range headers {
		if elem.Name == key {
			return elem.Value
		}
	}

	return ""
}

func FindEmailMessage(part []*gmail.MessagePart) (textMessage string) {
	for _, elem := range part {
		if elem.MimeType == "text/plain" {
			textMessage, err := utils.DecodeBase64(elem.Body.Data)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error when decoding: %+v", elem.Body.Data))
				return ""
			}

			return textMessage
		}
	}

	return ""
}

func (service *EmailService) ProcessWebhook(rawMessage string, prevHistoryId int) (historyId uint64, err error) {
	jsonString, err := utils.DecodeBase64(rawMessage)
	if err != nil {
		return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling DecodeBase64, error: %+v", err)
	}

	fmt.Println(fmt.Sprintf("===== json string %+v", jsonString))

	type WebhookData struct {
		EmailAddress string `json:"emailAddress"`
		HistoryId    uint64 `json:"historyId"`
	}

	var data WebhookData

	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling unmarshaling: %+v, error: %+v", jsonString, err)
	}

	fmt.Println(fmt.Sprintf("===== data %+v", data))

	// historyList, err := service.emailRepo.GetHistoryList(data.HistoryId)
	historyList, err := service.emailRepo.GetHistoryList(uint64(prevHistoryId))
	if err != nil {
		return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling GetHistoryList, history id: %+v, error: %+v", data.HistoryId, err)
	}

	fmt.Println(fmt.Sprintf("===== history list %+v", historyList))

	profile, err := service.emailRepo.GetProfile()
	if err != nil {
		return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling GetProfile, error: %+v", err)
	}
	fmt.Println(fmt.Sprintf("===== profile %+v", profile))

	for _, history := range historyList.History {
		if len(history.Messages) > 0 {
			message, err := service.emailRepo.GetMessageById(history.Messages[0].Id)
			if err != nil {
				return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling GetMessageById, message id: %+v, error: %+v", history.Messages[9].Id, err)
			}

			threadId := message.ThreadId
			subject := FindHeaders(message.Payload.Headers, "Subject")
			from := FindHeaders(message.Payload.Headers, "From")
			date := FindHeaders(message.Payload.Headers, "Date")
			emailMessage := FindEmailMessage(message.Payload.Parts)

			thread := entity.Thread{}
			interaction := &entity.Interaction{}

			ongoingInteraction, err := service.interactionRepo.GetOngoingInteractionByConversationId(threadId)
			if (ongoingInteraction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				existingThread, err := service.threadRepo.GetThreadByID(threadId)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newThread, err := service.threadRepo.InsertThread(entity.Thread{
						ID:        threadId,
						Subject:   subject,
						EmailDate: date,
						From:      from,
					})
					if err != nil {
						return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling InsertThread, error: %+v", err)
					}

					thread = newThread

				} else if err != nil {
					return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling GetThreadByID, error: %+v", err)

				} else {
					thread = existingThread
				}

				var reporterId uint
				var name string
				var email string

				name, email = utils.ParseFromHeader(from)
				if email == "" {
					email = from
				}

				existingReporter, err := service.reporterRepo.GetReporterByEmail(email)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newReporter := entity.Reporter{
						Name:  name,
						Email: email,
					}
					reporter, err := service.reporterRepo.CreateReporter(&newReporter)
					if err != nil {
						return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling InsertReporter, error: %+v", err)
					}

					reporterId = reporter.ID

				} else if err != nil {
					return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling GetReporterByEmail, error: %+v", err)

				} else {
					reporterId = existingReporter.ID
				}

				interaction = &entity.Interaction{
					ReporterId:      reporterId,
					Platform:        enum.EMAIL,
					InteractionType: enum.PESAN,
					Status:          enum.UNCLAIMED,
					ConversationId:  thread.ID,
					PlatformId:      profile.EmailAddress,
				}

				interaction, err = service.interactionRepo.CreateInteraction(interaction)
				if err != nil {
					return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling CreateInteraction, error: %+v", err)
				}

			} else if err != nil {
				return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling GetOngoingInteraction, error: %+v", err)

			} else {
				interaction = ongoingInteraction
			}

			m, err := service.messageRepo.GetMessageByMetaMessageId(history.Messages[0].Id)
			if m == nil { // No existing message so create new
				_, err = service.messageRepo.CreateMessage(&entity.Message{
					InteractionId: interaction.ID,
					Message:       emailMessage,
					SentBy:        enum.REPORTER,
					RecipientId:   profile.EmailAddress,
					MetaMessageId: history.Messages[0].Id,
				})
				if err != nil {
					return historyId, fmt.Errorf("[EmailService][ProcessWebhook] error when calling CreateMessage, error: %+v", err)
				}
			}

		}
	}

	// return data.HistoryId, nil
	return historyList.HistoryId, nil
}

func (service *EmailService) SendEmail(interactionId uint, message string) (messageId string, res *entity.Message, err error) {
	profile, err := service.emailRepo.GetProfile()
	if err != nil {
		return messageId, nil, fmt.Errorf("[EmailService][SendEmail] error when calling GetProfile, error: %+v", err)
	}

	interaction, err := service.interactionRepo.GetInteractionById(interactionId)
	if err != nil {
		return messageId, nil, fmt.Errorf("[EmailService][SendEmail] error when calling GetInteractionById, error: %+v", err)
	}

	thread, err := service.threadRepo.GetThreadByID(interaction.ConversationId)
	if err != nil {
		return messageId, nil, fmt.Errorf("[EmailService][SendEmail] error when calling GetThreadByID, error: %+v", err)
	}
	re := regexp.MustCompile(`<([^>]+)>`)

	// Finding submatches
	submatches := re.FindStringSubmatch(thread.From)

	// Extracting the string inside <>
	destination := ""
	if len(submatches) >= 2 {
		destination = submatches[1]
		fmt.Println("String inside <>:", destination)
	}

	latestMessage, err := service.messageRepo.GetLatestMessageofInteraction(interactionId)
	if err != nil {
		return messageId, nil, fmt.Errorf("[EmailService][SendEmail] error when calling GetLatestMessageofInteraction, error: %+v", err)
	}

	formattedMessage := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nReferences: %s\r\nIn-Reply-To: %s\r\n\r\n%s", profile.EmailAddress, destination, thread.Subject, latestMessage.MetaMessageId, latestMessage.MetaMessageId, message)

	gmailMessage := &gmail.Message{
		ThreadId: thread.ID,
		Raw:      base64.StdEncoding.EncodeToString([]byte(formattedMessage)),
	}

	gmailMessage, err = service.emailRepo.SendEmail(gmailMessage)
	if err != nil {
		return messageId, nil, fmt.Errorf("[EmailService][SendEmail] error when calling SendEmail, error: %+v", err)
	}

	res, err = service.messageRepo.CreateMessage(&entity.Message{
		InteractionId: interaction.ID,
		Message:       message,
		SentBy:        enum.AGENT,
		RecipientId:   profile.EmailAddress,
	})
	if err != nil {
		return messageId, nil, fmt.Errorf("[EmailService][SendEmail] error when calling CreateMessage, error: %+v", err)
	}

	return res.MetaMessageId, res, nil
}
