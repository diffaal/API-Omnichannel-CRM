package presentation

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/package/enum"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ID               uint       `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	InteractionId    uint       `json:"interaction_id"`
	SenderId         string     `json:"sender_id"`
	RecipientId      string     `json:"recipient_id"`
	MetaMessageId    string     `json:"mid"`
	Message          string     `json:"message"`
	MessageTimestamp time.Time  `json:"message_timestamp"`
	SentBy           string     `json:"sent_by"`
	IsRead           bool       `json:"is_read"`
	IsDeleted        bool       `json:"is_deleted"`
}

type MetaInteractionRequest struct {
	Object string       `json:"object"`
	Entry  []EntryField `json:"entry"`
}

type FacebookWebhookRequest struct {
	Object string               `json:"object"`
	Entry  []FacebookEntryField `json:"entry"`
}

type InstagramWebhookRequest struct {
	Object string                `json:"object"`
	Entry  []InstagramEntryField `json:"entry"`
}

type FacebookEntryField struct {
	PlatformId string                 `json:"id"`
	EntryTime  int64                  `json:"time"`
	Messaging  []MessagingField       `json:"messaging"`
	Changes    []FacebookChangesField `json:"changes"`
}

type InstagramEntryField struct {
	PlatformId string                  `json:"id"`
	EntryTime  int64                   `json:"time"`
	Messaging  []MessagingField        `json:"messaging"`
	Changes    []InstagramChangesField `json:"changes"`
}

type FacebookChangesField struct {
	Value     FacebookValueField `json:"value"`
	FieldType string             `json:"field"`
}

type FacebookValueField struct {
	Message   string `json:"message"`
	PostId    string `json:"post_id"`
	CommentId string `json:"comment_id"`
}

type InstagramChangesField struct {
	Value     InstagramValueField `json:"value"`
	FieldType string              `json:"field"`
}

type InstagramValueField struct {
	MediaId   string `json:"media_id"`
	CommentId string `json:"comment_id"`
}

type EntryField struct {
	PlatformId string           `json:"id"`
	EntryTime  int64            `json:"time"`
	Messaging  []MessagingField `json:"messaging"`
}

type MessagingField struct {
	ReporterId       IdField      `json:"sender"`
	PlatformId       IdField      `json:"recipient"`
	MessageTimestamp int64        `json:"timestamp"`
	Message          MessageField `json:"message"`
}

type IdField struct {
	Id string `json:"id"`
}

type MessageField struct {
	MessageId   string            `json:"mid"`
	MessageText string            `json:"text"`
	Attachments []AttachmentField `json:"attachments"`
}

type AttachmentField struct {
	AttachmentType    string                 `json:"type"`
	AttachmentPayload AttachmentPayloadField `json:"payload"`
}

type AttachmentPayloadField struct {
	AttachmentUrl string `json:"url"`
}

type ClaimInteractionRequest struct {
	InteractionId uint `json:"interaction_id"`
}

type DashboardInteractionList struct {
	InteractionId   uint      `json:"interaction_id"`
	PlatformId      string    `json:"platform_id"`
	ReporterId      uint      `json:"reporter_id"`
	ConversationId  string    `json:"conversation_id"`
	MentionMediaId  string    `json:"media_id"`
	AgentId         string    `json:"agent_id"`
	AgentName       string    `json:"agent_name"`
	Status          string    `json:"status"`
	Platform        string    `json:"platform"`
	InteractionType string    `json:"interaction_type"`
	Duration        time.Time `json:"duration"`
}

type MetaSendMessageRequest struct {
	InteractionId uint   `json:"interaction_id"`
	PlatformId    string `json:"platform_id,omitempty"`
	ReporterId    uint   `json:"reporter_id,omitempty"`
	Message       string `json:"message"`
	Platform      string `json:"platform"`
	SentBy        string `json:"sent_by"`
}

type MessengerSendMessageMetaRequest struct {
	Recipient     IdField              `json:"recipient"`
	MessagingType string               `json:"messaging_type"`
	Tag           string               `json:"tag"`
	Message       MessageSendMetaField `json:"message"`
}

type MessengerSendEmailRequest struct {
	InteractionId uint   `json:"interaction_id"`
	Message       string `json:"message"`
}

type MessageSendMetaField struct {
	MessageText string `json:"text"`
}

type MessengerSendMessageMetaResponse struct {
	ReporterId string `json:"recipient_id"`
	MessageId  string `json:"message_id"`
}

type InteractionLastMessage struct {
	Interaction entity.Interaction `json:"interaction"`
	LastMessage entity.Message     `json:"last_message"`
}

type WhatsappSendMessageMetaRequest struct {
	MessagingProduct string                   `json:"messaging_product"`
	RecipientType    string                   `json:"recipient_type"`
	Recipient        string                   `json:"to"`
	MessageType      string                   `json:"type"`
	Text             WhatsappMessageMetaField `json:"text"`
}

type WhatsappMessageMetaField struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type WhatsappSendMessageMetaResponse struct {
	MessagingProduct string                       `json:"messaging_product"`
	Contacts         []WhatsappContactMetaField   `json:"contacts"`
	Messages         []WhatsappMessageIdMetaField `json:"messages"`
}

type WhatsappContactMetaField struct {
	Input string `json:"input"`
	WaId  string `json:"wa_id"`
}

type WhatsappMessageIdMetaField struct {
	Id string `json:"id"`
}

type InteractionReporterData struct {
	Interaction entity.Interaction `json:"interaction"`
	Reporter    entity.Reporter    `json:"reporter"`
}

type SendInteractionDataCRMRequest struct {
	InteractionId       uint      `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           time.Time `json:"deleted_at"`
	ReporterName        string    `json:"name"`
	ReporterGender      string    `json:"gender"`
	ReporterPhoneNumber string    `json:"phoneNumber"`
	ReporterEmail       string    `json:"email"`
	ReporterAddress     string    `json:"address"`
	PlatformId          string    `json:"platformId"`
	ConversationId      string    `json:"conversationId"`
	ReporterId          uint      `json:"reporterId"`
	MentionMediaId      string    `json:"mentionMediaId"`
	MentionMediaUrl     string    `json:"mentionMediaUrl"`
	AgentId             string    `json:"agentId"`
	Status              string    `json:"status"`
	Platform            string    `json:"platform"`
	InteractionType     string    `json:"interactionType"`
	Latitude            string    `json:"latitude"`
	Longitude           string    `json:"longitude"`
	Duration            time.Time `json:"duration"`
}

type GmailInteractionRequest struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type GetMentionCommentDetailResp struct {
	MentionedComment MentionedCommentDetail `json:"mentioned_comment"`
	PlatformId       string                 `json:"id"`
}

type MentionedCommentDetail struct {
	Text      string           `json:"text"`
	Media     MediaDetailField `json:"media"`
	CommentId string           `json:"id"`
}

type MediaDetailField struct {
	Id        string `json:"id"`
	Permalink string `jsom:"permalink"`
}

type GetMentionMediaDetailResp struct {
	MentionedMedia MentionedMediaDetail `json:"mentioned_media"`
	PlatformId     string               `json:"id"`
}

type MentionedMediaDetail struct {
	Caption   string `json:"caption"`
	Permalink string `json:"permalink"`
	MediaId   string `json:"id"`
}

type CreateLiveChatInteractionRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type GetGeotagInformation struct {
	InterractionId uint   `json:"interraction_id" binding:"required"`
	Address        string `json:"address"`
	Lat            string `json:"lat"`
	Lon            string `json:"lon"`
}

type InteractionWithLatestMessage struct {
	InteractionId          uint      `json:"interaction_id"`
	InteractionCreatedAt   time.Time `json:"interaction_created_at"`
	PlatformId             string    `json:"platform_id"`
	ReporterId             uint      `json:"reporter_id"`
	ReporterName           string    `json:"reporter_name"`
	ConversationId         string    `json:"conversation_id"`
	MentionMediaId         string    `json:"media_id"`
	MentionMediaUrl        string    `json:"media_url"`
	AgentId                string    `json:"agent_id"`
	Status                 string    `json:"status"`
	Platform               string    `json:"platform"`
	InteractionType        string    `json:"interaction_type"`
	Latitude               string    `json:"Latitude"`
	Longitude              string    `json:"Longitude"`
	LatestMessageId        uint      `json:"latest_message_id"`
	LatestMessageCreatedAt time.Time `json:"latest_message_created_at"`
	SenderId               string    `json:"sender_id"`
	RecipientId            string    `json:"recipient_id"`
	MetaMessageId          string    `json:"mid"`
	Message                string    `json:"message"`
	AttachmentType         string    `json:"attachment_type"`
	AttachmentUrl          string    `json:"attachment_url"`
	SentBy                 string    `json:"sent_by"`
	IsRead                 bool      `json:"is_read"`
}

func (cir *ClaimInteractionRequest) ValidatePayload() map[string]string {
	errorMessage := make(map[string]string)

	if cir.InteractionId == 0 {
		errorMessage["errorStatus"] = enum.FIELD_REQUIRED_STATUS
		errorMessage["errorMessage"] = enum.FIELD_REQUIRED_MESSAGE
		return errorMessage
	}

	return errorMessage
}

func ParseGetListInteractionFilters(c *gin.Context) (map[string]interface{}, error) {
	filters := make(map[string]interface{})

	interactionIdsQuery := c.Query("interaction_ids")
	reporterIdsQuery := c.Query("reporter_ids")
	agentIdsQuery := c.Query("agent_ids")
	statusQuery := c.Query("status")
	platformsQuery := c.Query("platforms")
	interactionTypesQuery := c.Query("interaction_types")
	pageQuery := c.Query("page")
	pageSizeQuery := c.Query("pageSize")

	if interactionIdsQuery != "" {
		interactionIdsString := strings.Split(interactionIdsQuery, ",")
		var interactionIds []uint

		for _, v := range interactionIdsString {
			interactionId, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, err
			}
			interactionIdUint := uint(interactionId)
			interactionIds = append(interactionIds, interactionIdUint)
		}
		filters["interaction_ids"] = interactionIds
	}

	if reporterIdsQuery != "" {
		reporterIds := strings.Split(reporterIdsQuery, ",")
		filters["reporter_ids"] = reporterIds
	}

	if agentIdsQuery != "" {
		agentIds := strings.Split(agentIdsQuery, ",")
		filters["agent_ids"] = agentIds
	}

	if statusQuery != "" {
		status := strings.Split(statusQuery, ",")
		filters["status"] = status
	}

	if platformsQuery != "" {
		platforms := strings.Split(platformsQuery, ",")
		filters["platforms"] = platforms
	}

	if interactionTypesQuery != "" {
		interactionTypes := strings.Split(interactionTypesQuery, ",")
		filters["interaction_types"] = interactionTypes
	}

	if pageQuery != "" {
		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			return nil, err
		}
		filters["page"] = page
	}

	if pageSizeQuery != "" {
		pageSize, err := strconv.Atoi(pageSizeQuery)
		if err != nil {
			return nil, err
		}
		filters["pageSize"] = pageSize
	}

	return filters, nil
}
