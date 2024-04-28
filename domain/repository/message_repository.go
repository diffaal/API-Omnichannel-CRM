package repository

import (
	"Omnichannel-CRM/domain/entity"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

type IMessageRepository interface {
	CreateMessage(*entity.Message) (*entity.Message, error)
	GetMessageByMetaMessageId(string) (*entity.Message, error)
	GetMessagesofInteraction(uint) ([]entity.Message, error)
	GetLatestMessageofInteraction(uint) (*entity.Message, error)
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	messageRepo := MessageRepository{
		db: db,
	}

	return &messageRepo
}

func (mr *MessageRepository) CreateMessage(message *entity.Message) (*entity.Message, error) {
	err := mr.db.Create(&message).Error

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (mr *MessageRepository) GetMessageByMetaMessageId(metaMessageId string) (*entity.Message, error) {
	var message entity.Message

	err := mr.db.Where("meta_message_id = ?", metaMessageId).Take(&message).Error

	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (mr *MessageRepository) GetMessagesofInteraction(interactionId uint) ([]entity.Message, error) {
	var messages []entity.Message

	result := mr.db.Where("interaction_id = ?", interactionId).Order("created_at ASC, id ASC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return messages, nil
}

func (mr *MessageRepository) GetLatestMessageofInteraction(interactionId uint) (*entity.Message, error) {
	var message entity.Message

	result := mr.db.Where("interaction_id = ?", interactionId).Order("created_at DESC, id DESC").First(&message)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &message, nil
}
