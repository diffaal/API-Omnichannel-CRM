package repository

import (
	"Omnichannel-CRM/domain/entity"

	"gorm.io/gorm"
)

type ChannelAccountRepository struct {
	db *gorm.DB
}

type IChannelAccountRepository interface {
	CreateChannelAccount(*entity.ChannelAccount) (*entity.ChannelAccount, error)
	UpdateChannelAccount(uint, *entity.ChannelAccount) (*entity.ChannelAccount, error)
	GetChannelAccountList() ([]entity.ChannelAccount, error)
	GetChannelAccountById(uint) (*entity.ChannelAccount, error)
	DeleteChannelAccount(uint) error
}

func NewChannelAccountRepository(db *gorm.DB) *ChannelAccountRepository {
	channelAccountRepo := ChannelAccountRepository{
		db: db,
	}

	return &channelAccountRepo
}

func (car *ChannelAccountRepository) CreateChannelAccount(channelAccount *entity.ChannelAccount) (*entity.ChannelAccount, error) {
	err := car.db.Create(&channelAccount).Error

	if err != nil {
		return nil, err
	}

	return channelAccount, nil
}

func (car *ChannelAccountRepository) UpdateChannelAccount(id uint, newChannelAccount *entity.ChannelAccount) (*entity.ChannelAccount, error) {
	var currentChannelAccount entity.ChannelAccount

	err := car.db.Where("id = ?", id).First(&currentChannelAccount).Error
	if err != nil {
		return nil, err
	}

	if newChannelAccount.Name != "" {
		currentChannelAccount.Name = newChannelAccount.Name
	}

	if newChannelAccount.FaceboookPageId != "" {
		currentChannelAccount.FaceboookPageId = newChannelAccount.FaceboookPageId
	}

	if newChannelAccount.InstagramId != "" {
		currentChannelAccount.InstagramId = newChannelAccount.InstagramId
	}

	if newChannelAccount.WhatsappNumId != "" {
		currentChannelAccount.WhatsappNumId = newChannelAccount.WhatsappNumId
	}

	if newChannelAccount.WhatsappBusinessId != "" {
		currentChannelAccount.WhatsappBusinessId = newChannelAccount.WhatsappBusinessId
	}

	if newChannelAccount.FacebookAccessToken != "" {
		currentChannelAccount.FacebookAccessToken = newChannelAccount.FacebookAccessToken
	}

	if newChannelAccount.InstagramAccessToken != "" {
		currentChannelAccount.InstagramAccessToken = newChannelAccount.InstagramAccessToken
	}

	if newChannelAccount.WhatsappAccessToken != "" {
		currentChannelAccount.WhatsappAccessToken = newChannelAccount.WhatsappAccessToken
	}

	err = car.db.Save(&currentChannelAccount).Error
	if err != nil {
		return nil, err
	}

	return &currentChannelAccount, nil
}

func (car *ChannelAccountRepository) GetChannelAccountList() ([]entity.ChannelAccount, error) {
	var channelAccountList []entity.ChannelAccount

	result := car.db.Find(&channelAccountList)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return channelAccountList, nil
}

func (car *ChannelAccountRepository) GetChannelAccountById(channelAccountId uint) (*entity.ChannelAccount, error) {
	var channelAccount entity.ChannelAccount

	err := car.db.Where("id = ?", channelAccountId).Take(&channelAccount).Error

	if err != nil {
		return nil, err
	}
	return &channelAccount, nil
}

func (car *ChannelAccountRepository) DeleteChannelAccount(channelAccountId uint) error {
	var channelAccount entity.ChannelAccount

	err := car.db.Unscoped().Where("id = ?", channelAccountId).Delete(&channelAccount).Error

	if err != nil {
		return err
	}

	return nil
}
