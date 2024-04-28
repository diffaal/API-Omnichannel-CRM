package presentation

type ChannelAccountModel struct {
	Name                 string `json:"name"`
	FaceboookPageId      string `json:"facebook_page_id"`
	InstagramId          string `json:"instagram_id"`
	WhatsappNumId        string `json:"wa_num_id"`
	WhatsappBusinessId   string `json:"wa_business_id"`
	FacebookAccessToken  string `json:"facebook_access_token"`
	InstagramAccessToken string `json:"instagram_access_token"`
	WhatsappAccessToken  string `json:"whatsapp_access_token"`
}

type UpdateChannelAccountModel struct {
	ChannelAccountId     uint   `json:"channel_account_id"`
	Name                 string `json:"name"`
	FaceboookPageId      string `json:"facebook_page_id"`
	InstagramId          string `json:"instagram_id"`
	WhatsappNumId        string `json:"wa_num_id"`
	WhatsappBusinessId   string `json:"wa_business_id"`
	FacebookAccessToken  string `json:"facebook_access_token"`
	InstagramAccessToken string `json:"instagram_access_token"`
	WhatsappAccessToken  string `json:"whatsapp_access_token"`
}

type DeleteChannelAccountModel struct {
	ChannelAccountId uint `json:"channel_account_id"`
}
