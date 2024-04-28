package repository

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/presentation"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InteractionRepository struct {
	db *gorm.DB
}

type IinteractionRepository interface {
	CreateInteraction(*entity.Interaction) (*entity.Interaction, error)
	GetOngoingInteractionByReporterId(uint) (*entity.Interaction, error)
	GetOngoingInteractionByMentionMediaId(string) (*entity.Interaction, error)
	GetOngoingInteractionByConversationId(string) (*entity.Interaction, error)
	UpdateInteraction(uint, *entity.Interaction) (*entity.Interaction, error)
	GetInteractionList(map[string]interface{}, *entity.ChannelAccount) ([]entity.Interaction, int64, error)
	GetInteractionById(uint) (*entity.Interaction, error)
	GetInteractionsByAgentId(string, map[string]interface{}) ([]presentation.InteractionWithLatestMessage, error)
	GetActiveInteractionCount(string) (int64, error)
	GetInteractionHandledTodayCount(string) (int64, error)
	GetInteractionByConversationId(conversationid string) (*entity.Interaction, error)
}

func NewInteractionRepository(db *gorm.DB) *InteractionRepository {
	interactionRepo := InteractionRepository{
		db: db,
	}

	return &interactionRepo
}

func (ir *InteractionRepository) CreateInteraction(interaction *entity.Interaction) (*entity.Interaction, error) {
	err := ir.db.Create(&interaction).Error

	if err != nil {
		return nil, err
	}

	return interaction, nil
}

func (ir *InteractionRepository) GetOngoingInteractionByReporterId(reporterId uint) (*entity.Interaction, error) {
	var interaction entity.Interaction
	ongoingInteractionStatus := []string{
		enum.UNCLAIMED,
		enum.WAITING,
		enum.IN_PROGRESS,
		enum.ACTIVE,
		enum.INACTIVE,
		enum.MISSED,
		enum.UNPROCESSED,
		enum.PROCESSED,
	}

	result := ir.db.Where("reporter_id = ? AND status IN ?", reporterId, ongoingInteractionStatus).Take(&interaction)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &interaction, nil
}

func (ir *InteractionRepository) GetOngoingInteractionByMentionMediaId(mediaId string) (*entity.Interaction, error) {
	var interaction entity.Interaction
	ongoingInteractionStatus := []string{
		enum.UNCLAIMED,
		enum.WAITING,
		enum.IN_PROGRESS,
		enum.ACTIVE,
		enum.INACTIVE,
		enum.MISSED,
		enum.UNPROCESSED,
		enum.PROCESSED,
	}

	result := ir.db.Where("mention_media_id = ? AND status IN ?", mediaId, ongoingInteractionStatus).Take(&interaction)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &interaction, nil
}

func (ir *InteractionRepository) GetOngoingInteractionByConversationId(conversationId string) (*entity.Interaction, error) {
	var interaction entity.Interaction
	ongoingInteractionStatus := []string{
		enum.UNCLAIMED,
		enum.WAITING,
		enum.IN_PROGRESS,
		enum.ACTIVE,
		enum.INACTIVE,
		enum.MISSED,
		enum.UNPROCESSED,
		enum.PROCESSED,
	}

	result := ir.db.Where("conversation_id = ? AND status IN ?", conversationId, ongoingInteractionStatus).Take(&interaction)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &interaction, nil
}

func (ir *InteractionRepository) UpdateInteraction(interactionId uint, newInteraction *entity.Interaction) (*entity.Interaction, error) {
	var currentInteraction entity.Interaction

	err := ir.db.Where("id = ?", interactionId).First(&currentInteraction).Error
	if err != nil {
		return nil, err
	}

	if newInteraction.AgentId != "" {
		currentInteraction.AgentId = newInteraction.AgentId
	}

	if newInteraction.Status != "" {
		currentInteraction.Status = newInteraction.Status
	}

	if newInteraction.Latitude != "" {
		currentInteraction.Latitude = newInteraction.Latitude
	}

	if newInteraction.Longitude != "" {
		currentInteraction.Longitude = newInteraction.Longitude
	}

	err = ir.db.Save(&currentInteraction).Error
	if err != nil {
		return nil, err
	}

	return &currentInteraction, nil
}

func (ir *InteractionRepository) GetInteractionList(filters map[string]interface{}, channelAccount *entity.ChannelAccount) ([]entity.Interaction, int64, error) {
	var interactionList []entity.Interaction
	var count int64
	queryDB := ir.db
	initialCondition := "platform = 'EMAIL'"

	if channelAccount.ID != 0 {
		if channelAccount.FaceboookPageId != "" {
			initialCondition = initialCondition + fmt.Sprintf(" OR platform_id = '%s'", channelAccount.FaceboookPageId)
		}
		if channelAccount.InstagramId != "" {
			initialCondition = initialCondition + fmt.Sprintf(" OR platform_id = '%s'", channelAccount.InstagramId)
		}
		if channelAccount.WhatsappBusinessId != "" {
			initialCondition = initialCondition + fmt.Sprintf(" OR platform_id = '%s'", channelAccount.WhatsappBusinessId)
		}
		if channelAccount.IsLiveChatActive == true {
			initialCondition = initialCondition + " OR platform = 'LIVE_CHAT'"
		}
	}

	queryDB = queryDB.Or(initialCondition)

	if filters["interaction_ids"] != nil {
		queryDB = queryDB.Where("id IN ?", filters["interaction_ids"])
	}
	if filters["reporter_ids"] != nil {
		queryDB = queryDB.Where("reporter_id IN ?", filters["reporter_ids"])
	}
	if filters["agent_ids"] != nil {
		queryDB = queryDB.Where("agent_id IN ?", filters["agent_ids"])
	}
	if filters["status"] != nil {
		queryDB = queryDB.Where("status IN ?", filters["status"])
	}
	if filters["platforms"] != nil {
		queryDB = queryDB.Where("platform IN ?", filters["platforms"])
	}
	if filters["interaction_types"] != nil {
		queryDB = queryDB.Where("interaction_type IN ?", filters["interaction_types"])
	}

	err := queryDB.Model(&entity.Interaction{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	if filters["page"] != nil && filters["pageSize"] != nil {
		offset := (filters["page"].(int) - 1) * filters["pageSize"].(int)
		limit := filters["pageSize"].(int)
		queryDB = queryDB.Offset(offset).Limit(limit)
	}

	result := queryDB.Order("created_at DESC, id DESC").Find(&interactionList)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, 0, nil
	}

	return interactionList, count, nil
}

func (ir *InteractionRepository) GetInteractionById(interactionId uint) (*entity.Interaction, error) {
	var interaction entity.Interaction

	result := ir.db.Where("id = ?", interactionId).Take(&interaction)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &interaction, nil
}

func (ir *InteractionRepository) GetInteractionsByAgentId(agentId string, filters map[string]interface{}) ([]presentation.InteractionWithLatestMessage, error) {
	var interactions []presentation.InteractionWithLatestMessage
	var interaction presentation.InteractionWithLatestMessage

	statement := fmt.Sprintf(`SELECT
		interactions.id AS interaction_id, 
		interactions.created_at AS interaction_created_at,
		interactions.platform_id,
		interactions.reporter_id,
		interactions.conversation_id,
		interactions.mention_media_id,
		interactions.agent_id,
		interactions.status,
		interactions.platform,
		interactions.interaction_type,
		interactions.mention_media_url,
		interactions.latitude,
		interactions.longitude,
		reporters.name AS reporter_name,
		latest_message.id AS latest_message_id,
		latest_message.created_at AS latest_message_created_at,
		latest_message.sender_id,
		latest_message.recipient_id,
		latest_message.meta_message_id,
		latest_message.message,
		latest_message.sent_by,
		latest_message.is_read,
		latest_message.attachment_type,
		latest_message.attachment_url
	FROM interactions
	INNER JOIN (
		SELECT interaction_id, MAX(created_at) AS latest_timestamp
		FROM messages
		GROUP BY interaction_id
	) AS latest_messages
	ON interactions.id = latest_messages.interaction_id
	INNER JOIN messages AS latest_message
	ON latest_messages.interaction_id = latest_message.interaction_id
	AND latest_messages.latest_timestamp = latest_message.created_at
	INNER JOIN reporters
	ON interactions.reporter_id = reporters.id
	WHERE interactions.agent_id = '%s'`, agentId)

	if filters["status"] != nil {
		formattedStatus := "('" + strings.Join(filters["status"].([]string), "', '") + "')"
		statement = statement + " AND status IN " + formattedStatus
	}
	if filters["platforms"] != nil {
		formattedPlatforms := "('" + strings.Join(filters["platforms"].([]string), "', '") + "')"
		statement = statement + " AND platform IN " + formattedPlatforms
	}
	if filters["interaction_types"] != nil {
		formattedTypes := "('" + strings.Join(filters["interaction_types"].([]string), "', '") + "')"
		statement = statement + " AND interaction_type IN " + formattedTypes
	}

	statement = statement + " ORDER BY latest_message.created_at DESC"

	rows, err := ir.db.Raw(statement).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		ir.db.ScanRows(rows, &interaction)
		interactions = append(interactions, interaction)
	}
	return interactions, nil
}

func (ir *InteractionRepository) GetActiveInteractionCount(agentId string) (int64, error) {
	var count int64
	status := []string{enum.IN_PROGRESS, enum.ACTIVE}
	err := ir.db.Model(&entity.Interaction{}).Where("agent_id = ? AND status IN ? ", agentId, status).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ir *InteractionRepository) GetInteractionHandledTodayCount(agentId string) (int64, error) {
	var count int64
	currentDateTime := time.Now()
	todayDate := fmt.Sprintf(currentDateTime.Format("2006-01-02"))
	startDate := fmt.Sprintf("%s 00:00:00", todayDate)
	endDate := fmt.Sprintf("%s 23:59:59", todayDate)

	err := ir.db.Model(&entity.Interaction{}).Where("agent_id = ? AND updated_at BETWEEN ? AND ? ", agentId, startDate, endDate).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ir *InteractionRepository) GetInteractionByConversationId(conversationid string) (*entity.Interaction, error) {
	var interaction entity.Interaction

	result := ir.db.Where("conversation_id = ?", conversationid).First(&interaction)
	if result.Error != nil {
		return &interaction, result.Error
	}
	if result.RowsAffected == 0 {
		return &interaction, nil
	}

	return &interaction, nil
}
