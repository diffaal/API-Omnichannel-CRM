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
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type MetaWebhookService struct {
	interactionRepo repository.IinteractionRepository
	messageRepo     repository.IMessageRepository
	reporterRepo    repository.IReporterRepository
}

type IMetaWebhookService interface {
	FacebookInteractionService(*presentation.FacebookWebhookRequest) (map[string]interface{}, []entity.Message, error)
	InstagramInteractionService(*presentation.InstagramWebhookRequest) (map[string]interface{}, []entity.Message, error)
	WhatsappMessageInteractionService(*presentation.WhatsappInteractionRequest) (map[string]interface{}, []entity.Message, error)
	WebsocketSendService(messages []entity.Message) error
}

func NewMetaWebhookService(interactionRepo repository.IinteractionRepository, messageRepo repository.IMessageRepository, reporterRepo repository.IReporterRepository) *MetaWebhookService {
	metaWebhookService := MetaWebhookService{
		interactionRepo: interactionRepo,
		messageRepo:     messageRepo,
		reporterRepo:    reporterRepo,
	}
	return &metaWebhookService
}

type MessageToSend struct {
	Action  string               `json:"action"`
	Message presentation.Message `json:"message"`
}

func (mws *MetaWebhookService) WebsocketSendService(messages []entity.Message) error {
	config.GetConfig()

	host := viper.GetString("Websocket.Host")
	channel := "ws"

	for _, message := range messages {
		query := fmt.Sprintf("user_id=%s&room_id=%d", message.SenderId, message.InteractionId)
		client, err := NewWebSocketClient(host, channel, message.SenderId, query)
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
		client.Write(send)
	}

	return nil
}

func (mws *MetaWebhookService) WhatsappMessageInteractionService(mir *presentation.WhatsappInteractionRequest) (map[string]interface{}, []entity.Message, error) {
	result := make(map[string]interface{})
	var resStructList []entity.Message
	var resStruct entity.Message
	var messageIds []string
	var interactionIds []uint
	var reporterIds []uint
	var interactionId uint
	var reporterId uint
	var platform string

	resStruct.SentBy = enum.REPORTER

	entries := mir.Entry

	platform = enum.WA

	for _, v := range entries {
		platformId := v.PlatformId
		metaReporterId := v.Changes[0].Value.Contacts[0].WaId
		messageId := v.Changes[0].Value.Messages[0].ID
		messageText := v.Changes[0].Value.Messages[0].Text.Body
		messageType := v.Changes[0].Value.Messages[0].Type
		messageTimestampString := v.Changes[0].Value.Messages[0].Timestamp
		timestampInt, _ := strconv.ParseInt(messageTimestampString, 10, 64)
		messageTimestamp := time.Unix(timestampInt, 0)

		existingReporter, err := mws.reporterRepo.GetReporterByMetaReporterId(metaReporterId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newReporter := entity.Reporter{
				MetaReporterId: metaReporterId,
			}
			reporter, err := mws.reporterRepo.CreateReporter(&newReporter)
			if err != nil {
				return nil, nil, err
			}

			reporterId = reporter.ID
			reporterIds = append(reporterIds, reporter.ID)

		} else if err != nil {
			return nil, nil, err

		} else {
			reporterId = existingReporter.ID
			reporterIds = append(reporterIds, existingReporter.ID)
		}

		ongoingInteraction, err := mws.interactionRepo.GetOngoingInteractionByReporterId(reporterId)
		if (ongoingInteraction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
			newInteraction := entity.Interaction{
				PlatformId:      platformId,
				ReporterId:      reporterId,
				Status:          enum.UNCLAIMED,
				Platform:        platform,
				InteractionType: enum.PESAN,
			}

			interaction, err := mws.interactionRepo.CreateInteraction(&newInteraction)
			if err != nil {
				return nil, resStructList, err
			}

			interactionId = interaction.ID
			interactionIds = append(interactionIds, interactionId)

		} else if err != nil {
			return nil, resStructList, err

		} else {
			interactionId = ongoingInteraction.ID
			interactionIds = append(interactionIds, interactionId)

		}

		existingMessage, err := mws.messageRepo.GetMessageByMetaMessageId(messageId)

		if existingMessage != nil {
			resStruct = *existingMessage
		}

		if (existingMessage == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
			newMessage := entity.Message{
				InteractionId:    interactionId,
				SenderId:         metaReporterId,
				RecipientId:      platformId,
				MetaMessageId:    messageId,
				Message:          messageText,
				MessageTimestamp: messageTimestamp,
				SentBy:           enum.REPORTER,
				IsRead:           false,
			}

			if messageType == "location" {
				newMessage.AttachmentType = enum.LOCATION
				location, _ := json.Marshal(v.Changes[0].Value.Messages[0].Location)
				newMessage.Message = string(location)
			}

			if messageType == "image" || messageType == "video" {
				var attachmentId string
				var attachmentCaption string

				if messageType == "image" {
					newMessage.AttachmentType = enum.IMAGE
					attachmentId = v.Changes[0].Value.Messages[0].Image.Id
					attachmentCaption = v.Changes[0].Value.Messages[0].Image.Caption
				} else if messageType == "video" {
					newMessage.AttachmentType = enum.VIDEO
					attachmentId = v.Changes[0].Value.Messages[0].Video.Id
					attachmentCaption = v.Changes[0].Value.Messages[0].Video.Caption
				}

				if attachmentCaption != "" {
					newMessage.Message = attachmentCaption
				}

				path := fmt.Sprintf("/%s/%s", viper.GetString("Meta.API_VERSION"), attachmentId)
				accessToken := viper.GetString("Meta.WA_ACCESS_TOKEN")

				reqUrl := url.URL{
					Scheme: "https",
					Host:   "graph.facebook.com",
					Path:   path,
				}

				response, err := request.GetRequest(reqUrl, accessToken)
				if err != nil {
					return nil, nil, err
				}

				defer response.Body.Close()

				if response.StatusCode == http.StatusOK {
					attachDetailData := &presentation.WhatsappAttachmentDetailResp{}
					err = json.NewDecoder(response.Body).Decode(attachDetailData)
					if err != nil {
						return nil, nil, err
					}

					attachmentUrl := attachDetailData.Url
					newMessage.AttachmentUrl = attachmentUrl

				} else {
					bodyBytes, err := io.ReadAll(response.Body)
					if err != nil {
						return nil, nil, err
					}
					bodyString := string(bodyBytes)
					logger.Info(response.Status)
					logger.Info(bodyString)
					return nil, nil, errors.New("Meta response not 200")
				}
			}

			message, err := mws.messageRepo.CreateMessage(&newMessage)
			if err != nil {
				return nil, resStructList, err
			}

			resStruct = *message
			messageIds = append(messageIds, message.MetaMessageId)

		} else if err != nil {
			return nil, resStructList, err

		} else {
			messageIds = append(messageIds, existingMessage.MetaMessageId)
		}
		resStructList = append(resStructList, resStruct)
	}

	result["messageIds"] = messageIds
	result["interactionIds"] = interactionIds
	result["reporterIds"] = reporterIds

	return result, resStructList, nil
}

func (mws *MetaWebhookService) FacebookInteractionService(fwr *presentation.FacebookWebhookRequest) (map[string]interface{}, []entity.Message, error) {
	result := make(map[string]interface{})
	var resStructList []entity.Message
	var resStruct entity.Message
	var messageIds []string
	var interactionIds []uint
	var reporterIds []string
	var interactionId uint
	var reporterId uint

	entries := fwr.Entry

	for _, v := range entries {
		platformId := v.PlatformId
		messaging := v.Messaging
		changes := v.Changes

		if len(messaging) > 0 {
			metaReporterId := v.Messaging[0].ReporterId.Id
			messageId := v.Messaging[0].Message.MessageId
			messageText := v.Messaging[0].Message.MessageText
			attachments := v.Messaging[0].Message.Attachments

			if platformId == metaReporterId {
				continue
			}

			existingReporter, err := mws.reporterRepo.GetReporterByMetaReporterId(metaReporterId)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newReporter := entity.Reporter{
					MetaReporterId: metaReporterId,
				}
				reporter, err := mws.reporterRepo.CreateReporter(&newReporter)
				if err != nil {
					return nil, nil, err
				}

				reporterId = reporter.ID
				reporterIds = append(reporterIds, reporter.MetaReporterId)

			} else if err != nil {
				return nil, nil, err

			} else {
				reporterId = existingReporter.ID
				reporterIds = append(reporterIds, existingReporter.MetaReporterId)
			}

			ongoingInteraction, err := mws.interactionRepo.GetOngoingInteractionByReporterId(reporterId)
			if (ongoingInteraction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newInteraction := entity.Interaction{
					PlatformId:      platformId,
					ReporterId:      reporterId,
					Status:          enum.UNCLAIMED,
					Platform:        enum.FACEBOOK,
					InteractionType: enum.PESAN,
				}

				interaction, err := mws.interactionRepo.CreateInteraction(&newInteraction)
				if err != nil {
					return nil, nil, err
				}

				interactionId = interaction.ID
				interactionIds = append(interactionIds, interactionId)

			} else if err != nil {
				return nil, nil, err

			} else {
				interactionId = ongoingInteraction.ID
				interactionIds = append(interactionIds, interactionId)

			}

			existingMessage, err := mws.messageRepo.GetMessageByMetaMessageId(messageId)
			if existingMessage != nil {
				resStruct = *existingMessage
			}
			if (existingMessage == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newMessage := entity.Message{
					InteractionId: interactionId,
					SenderId:      metaReporterId,
					RecipientId:   platformId,
					MetaMessageId: messageId,
					Message:       messageText,
					SentBy:        enum.REPORTER,
					IsRead:        false,
				}

				if len(attachments) > 0 {
					attachmentType := attachments[0].AttachmentType
					attachmentUrl := attachments[0].AttachmentPayload.AttachmentUrl

					if attachmentType == "image" || attachmentType == "video" {
						if attachmentType == "image" {
							newMessage.AttachmentType = enum.IMAGE
						} else if attachmentType == "video" {
							newMessage.AttachmentType = enum.VIDEO
						}
						newMessage.AttachmentUrl = attachmentUrl
					}
				}

				message, err := mws.messageRepo.CreateMessage(&newMessage)
				if err != nil {
					return nil, nil, err
				}
				resStruct = *message

				messageIds = append(messageIds, message.MetaMessageId)

			} else if err != nil {
				return nil, nil, err

			} else {
				messageIds = append(messageIds, existingMessage.MetaMessageId)
			}
			resStructList = append(resStructList, resStruct)

		} else if len(changes) > 0 {
			postId := v.Changes[0].Value.PostId
			commentId := v.Changes[0].Value.CommentId
			mentionMessage := v.Changes[0].Value.Message
			mediaUrl := fmt.Sprintf("https://www.facebook.com/%s", postId)
			var mentionMessageId string

			if commentId != "" {
				mentionMessageId = commentId
			} else {
				mentionMessageId = postId
			}

			ongoingInteraction, err := mws.interactionRepo.GetOngoingInteractionByMentionMediaId(postId)
			if (ongoingInteraction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newInteraction := entity.Interaction{
					PlatformId:      platformId,
					MentionMediaId:  postId,
					MentionMediaUrl: mediaUrl,
					Status:          enum.UNCLAIMED,
					Platform:        enum.FACEBOOK,
					InteractionType: enum.MENTION,
				}

				interaction, err := mws.interactionRepo.CreateInteraction(&newInteraction)
				if err != nil {
					return nil, nil, err
				}

				interactionId = interaction.ID
				interactionIds = append(interactionIds, interactionId)

			} else if err != nil {
				return nil, nil, err

			} else {
				interactionId = ongoingInteraction.ID
				interactionIds = append(interactionIds, interactionId)

			}

			existingMessage, err := mws.messageRepo.GetMessageByMetaMessageId(mentionMessageId)
			if existingMessage != nil {
				resStruct = *existingMessage
			}
			if (existingMessage == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newMessage := entity.Message{
					InteractionId: interactionId,
					RecipientId:   platformId,
					MetaMessageId: mentionMessageId,
					Message:       mentionMessage,
					SentBy:        enum.REPORTER,
					IsRead:        false,
				}

				message, err := mws.messageRepo.CreateMessage(&newMessage)
				if err != nil {
					return nil, nil, err
				}

				resStruct = *message
				messageIds = append(messageIds, message.MetaMessageId)

			} else if err != nil {
				return nil, nil, err

			} else {
				messageIds = append(messageIds, existingMessage.MetaMessageId)
			}
			resStructList = append(resStructList, resStruct)

		} else {
			continue
		}
	}

	result["messageIds"] = messageIds
	result["interactionIds"] = interactionIds
	result["reporterIds"] = reporterIds

	return result, resStructList, nil
}

func (mws *MetaWebhookService) InstagramInteractionService(iwr *presentation.InstagramWebhookRequest) (map[string]interface{}, []entity.Message, error) {
	result := make(map[string]interface{})
	var resStructList []entity.Message
	var resStruct entity.Message
	var messageIds []string
	var interactionIds []uint
	var reporterIds []string
	var interactionId uint
	var reporterId uint

	entries := iwr.Entry

	for _, v := range entries {
		platformId := v.PlatformId
		messaging := v.Messaging
		changes := v.Changes

		if len(messaging) > 0 {
			metaReporterId := v.Messaging[0].ReporterId.Id
			messageId := v.Messaging[0].Message.MessageId
			messageText := v.Messaging[0].Message.MessageText
			attachments := v.Messaging[0].Message.Attachments

			if platformId == metaReporterId {
				continue
			}

			existingReporter, err := mws.reporterRepo.GetReporterByMetaReporterId(metaReporterId)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newReporter := entity.Reporter{
					MetaReporterId: metaReporterId,
				}
				reporter, err := mws.reporterRepo.CreateReporter(&newReporter)
				if err != nil {
					return nil, nil, err
				}

				reporterId = reporter.ID
				reporterIds = append(reporterIds, reporter.MetaReporterId)

			} else if err != nil {
				return nil, nil, err

			} else {
				reporterId = existingReporter.ID
				reporterIds = append(reporterIds, existingReporter.MetaReporterId)
			}

			ongoingInteraction, err := mws.interactionRepo.GetOngoingInteractionByReporterId(reporterId)
			if (ongoingInteraction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newInteraction := entity.Interaction{
					PlatformId:      platformId,
					ReporterId:      reporterId,
					Status:          enum.UNCLAIMED,
					Platform:        enum.IG,
					InteractionType: enum.PESAN,
				}

				interaction, err := mws.interactionRepo.CreateInteraction(&newInteraction)
				if err != nil {
					return nil, nil, err
				}

				interactionId = interaction.ID
				interactionIds = append(interactionIds, interactionId)

			} else if err != nil {
				return nil, nil, err

			} else {
				interactionId = ongoingInteraction.ID
				interactionIds = append(interactionIds, interactionId)

			}

			existingMessage, err := mws.messageRepo.GetMessageByMetaMessageId(messageId)
			if existingMessage != nil {
				resStruct = *existingMessage
			}
			if (existingMessage == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newMessage := entity.Message{
					InteractionId: interactionId,
					SenderId:      metaReporterId,
					RecipientId:   platformId,
					MetaMessageId: messageId,
					Message:       messageText,
					SentBy:        enum.REPORTER,
					IsRead:        false,
				}

				if len(attachments) > 0 {
					attachmentType := attachments[0].AttachmentType
					attachmentUrl := attachments[0].AttachmentPayload.AttachmentUrl

					if attachmentType == "image" || attachmentType == "video" {
						if attachmentType == "image" {
							newMessage.AttachmentType = enum.IMAGE
						} else if attachmentType == "video" {
							newMessage.AttachmentType = enum.VIDEO
						}
						newMessage.AttachmentUrl = attachmentUrl
					}
				}

				message, err := mws.messageRepo.CreateMessage(&newMessage)
				if err != nil {
					return nil, nil, err
				}
				resStruct = *message

				messageIds = append(messageIds, message.MetaMessageId)

			} else if err != nil {
				return nil, nil, err

			} else {
				messageIds = append(messageIds, existingMessage.MetaMessageId)
			}
			resStructList = append(resStructList, resStruct)

		} else if len(changes) > 0 {
			mediaId := v.Changes[0].Value.MediaId
			commentId := v.Changes[0].Value.CommentId
			var mediaUrl string
			var mentionMessage string
			var mentionMessageId string

			if commentId != "" {
				path := fmt.Sprintf("/%s/%s", viper.GetString("Meta.API_VERSION"), platformId)
				params := url.Values{}
				params.Add("fields", fmt.Sprintf("mentioned_comment.comment_id(%s){text,media{id,permalink}}", commentId))
				params.Add("access_token", viper.GetString("Meta.IG_ACCESS_TOKEN"))

				reqUrl := url.URL{
					Scheme:   "https",
					Host:     "graph.facebook.com",
					Path:     path,
					RawQuery: params.Encode(),
				}

				response, err := request.GetRequest(reqUrl, "")
				if err != nil {
					return nil, nil, err
				}

				defer response.Body.Close()

				if response.StatusCode == http.StatusOK {
					mentionCommentData := &presentation.GetMentionCommentDetailResp{}
					err = json.NewDecoder(response.Body).Decode(mentionCommentData)
					if err != nil {
						return nil, nil, err
					}

					mediaUrl = mentionCommentData.MentionedComment.Media.Permalink
					mentionMessage = mentionCommentData.MentionedComment.Text
					mentionMessageId = commentId

				} else {
					bodyBytes, err := io.ReadAll(response.Body)
					if err != nil {
						return nil, nil, err
					}
					bodyString := string(bodyBytes)
					logger.Info(response.Status)
					logger.Info(bodyString)
					return nil, nil, errors.New("Meta response not 200")
				}

			} else {
				path := fmt.Sprintf("/%s/%s", viper.GetString("Meta.API_VERSION"), platformId)
				params := url.Values{}
				params.Add("fields", fmt.Sprintf("mentioned_media.media_id(%s){caption,permalink}", mediaId))
				params.Add("access_token", viper.GetString("Meta.IG_ACCESS_TOKEN"))

				reqUrl := url.URL{
					Scheme:   "https",
					Host:     "graph.facebook.com",
					Path:     path,
					RawQuery: params.Encode(),
				}

				response, err := request.GetRequest(reqUrl, "")
				if err != nil {
					return nil, nil, err
				}

				defer response.Body.Close()

				if response.StatusCode == http.StatusOK {
					mentionMediaData := &presentation.GetMentionMediaDetailResp{}
					err = json.NewDecoder(response.Body).Decode(mentionMediaData)
					if err != nil {
						return nil, nil, err
					}

					mediaUrl = mentionMediaData.MentionedMedia.Permalink
					mentionMessage = mentionMediaData.MentionedMedia.Caption
					mentionMessageId = mediaId

				} else {
					bodyBytes, err := io.ReadAll(response.Body)
					if err != nil {
						return nil, nil, err
					}
					bodyString := string(bodyBytes)
					logger.Info(response.Status)
					logger.Info(bodyString)
					return nil, nil, errors.New("Meta response not 200")
				}
			}

			ongoingInteraction, err := mws.interactionRepo.GetOngoingInteractionByMentionMediaId(mediaId)
			if (ongoingInteraction == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newInteraction := entity.Interaction{
					PlatformId:      platformId,
					MentionMediaId:  mediaId,
					MentionMediaUrl: mediaUrl,
					Status:          enum.UNCLAIMED,
					Platform:        enum.IG,
					InteractionType: enum.MENTION,
				}

				interaction, err := mws.interactionRepo.CreateInteraction(&newInteraction)
				if err != nil {
					return nil, nil, err
				}

				interactionId = interaction.ID
				interactionIds = append(interactionIds, interactionId)

			} else if err != nil {
				return nil, nil, err

			} else {
				interactionId = ongoingInteraction.ID
				interactionIds = append(interactionIds, interactionId)

			}

			existingMessage, err := mws.messageRepo.GetMessageByMetaMessageId(mentionMessageId)
			if existingMessage != nil {
				resStruct = *existingMessage
			}
			if (existingMessage == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
				newMessage := entity.Message{
					InteractionId: interactionId,
					RecipientId:   platformId,
					MetaMessageId: mentionMessageId,
					Message:       mentionMessage,
					SentBy:        enum.REPORTER,
					IsRead:        false,
				}

				message, err := mws.messageRepo.CreateMessage(&newMessage)
				if err != nil {
					return nil, nil, err
				}
				resStruct = *message

				messageIds = append(messageIds, message.MetaMessageId)

			} else if err != nil {
				return nil, nil, err

			} else {
				messageIds = append(messageIds, existingMessage.MetaMessageId)
			}
			resStructList = append(resStructList, resStruct)
		} else {
			continue
		}
	}

	result["messageIds"] = messageIds
	result["interactionIds"] = interactionIds
	result["reporterIds"] = reporterIds

	return result, resStructList, nil
}
