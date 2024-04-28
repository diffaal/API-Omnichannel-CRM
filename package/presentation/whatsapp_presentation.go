package presentation

import "time"

type WhatsappInteractionRequest struct {
	Object string               `json:"object"`
	Entry  []WhatsappEntryField `json:"entry"`
}

type WhatsappEntryField struct {
	PlatformId string            `json:"id"`
	Changes    []WhatsappChanges `json:"changes"`
}

type WhatsappChanges struct {
	Value WhatsappValue `json:"value"`
	Field string        `json:"field"`
}

type WhatsappValue struct {
	MessagingProduct string             `json:"messaging_product"`
	Metadata         WhatsappMetadata   `json:"metadata"`
	Contacts         []WhatsappContacts `json:"contacts"`
	Errors           []WhatsappErrors   `json:"errors"`
	Messages         []WhatsappMessages `json:"messages"`
	Statuses         interface{}        `json:"statuses"`
}

type WhatsappMessages struct {
	ID        string           `json:"id"`
	Text      WhatsappText     `json:"text"`
	Image     WhatsappImage    `json:"image"`
	Video     WhatsappVideo    `json:"video"`
	Location  WhatsappLocation `json:"location"`
	Timestamp string           `json:"timestamp"`
	Type      string           `json:"type"`
}

type WhatsappText struct {
	Body string `json:"body"`
}

type WhatsappLocation struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type WhatsappImage struct {
	Caption  string `json:"caption"`
	Sha256   string `json:"sha256"`
	Id       string `json:"id"`
	MimeType string `json:"mime_type"`
}

type WhatsappVideo struct {
	Caption  string `json:"caption"`
	Filename string `json:"filename"`
	Sha256   string `json:"sha256"`
	Id       string `json:"id"`
	MimeType string `json:"mime_type"`
}
type WhatsappStatuses struct {
	ID          string         `json:"id"`
	RecipientId string         `json:"recipient_id"`
	Status      string         `json:"status"`
	Errors      WhatsappErrors `json:"errors"`
	Timestamp   time.Time      `json:"timestamp"`
}

type WhatsappConversation struct {
	ID                  string         `json:"id"`
	Origin              WhatsappOrigin `json:"origin"`
	ExpirationTimestamp string         `json:"expiration_timestamp"`
}

type WhatsappOrigin struct {
	Type string `json:"type"`
}

type WhatsappErrors struct {
	Code      int               `json:"code"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	ErrorData WhatsappErrorData `json:"error_data"`
}

type WhatsappErrorData struct {
	Details string `json:"details"`
}

type WhatsappContacts struct {
	WaId    string            `json:"wa_id"`
	Profile map[string]string `json:"profile"`
}

type WhatsappMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type WebsocketMessage struct {
	InteractionId uint   `json:"interaction_id"`
	PlatformId    string `json:"platform_id"`
	ReporterId    uint   `json:"reporter_id"`
	Message       string `json:"message"`
	Platform      string `json:"platform"`
}
type WhatsappAttachmentDetailResp struct {
	MessagingProduct string `json:"messaging_product"`
	Url              string `json:"url"`
	Sha256           string `json:"sha256"`
	Id               string `json:"id"`
	MimeType         string `json:"mime_type"`
}
