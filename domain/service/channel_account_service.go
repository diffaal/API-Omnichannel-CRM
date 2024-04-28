package service

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/presentation"

	"errors"

	"gorm.io/gorm"
)

type ChannelAccountService struct {
	channelAccountRepo repository.IChannelAccountRepository
}

type IChannelAccountService interface {
	CreateChannelAccount(*presentation.ChannelAccountModel) (map[string]interface{}, error)
	GetChannelAccountList() (map[string]interface{}, error)
	GetChannelAccountById(uint) (map[string]interface{}, error)
	UpdateChannelAccount(*presentation.UpdateChannelAccountModel) (map[string]interface{}, error)
	DeleteChannelAccountById(*presentation.DeleteChannelAccountModel) (map[string]interface{}, error)
}

func NewChannelAccountService(channelAccountRepo repository.IChannelAccountRepository) *ChannelAccountService {
	channelAccountService := ChannelAccountService{
		channelAccountRepo: channelAccountRepo,
	}
	return &channelAccountService
}

func (cas *ChannelAccountService) CreateChannelAccount(cam *presentation.ChannelAccountModel) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	newChannelAccount := entity.ChannelAccount{
		Name:                 cam.Name,
		FaceboookPageId:      cam.FaceboookPageId,
		InstagramId:          cam.InstagramId,
		WhatsappNumId:        cam.WhatsappNumId,
		WhatsappBusinessId:   cam.WhatsappBusinessId,
		FacebookAccessToken:  cam.FacebookAccessToken,
		InstagramAccessToken: cam.InstagramAccessToken,
		WhatsappAccessToken:  cam.WhatsappAccessToken,
	}

	channelAccount, err := cas.channelAccountRepo.CreateChannelAccount(&newChannelAccount)
	if err != nil {
		return nil, err
	}

	result["channel_account"] = channelAccount

	return result, nil
}

func (cas *ChannelAccountService) GetChannelAccountList() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	channelAccountList, err := cas.channelAccountRepo.GetChannelAccountList()
	if (channelAccountList == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	result["channel_account_list"] = channelAccountList

	return result, nil
}

func (cas *ChannelAccountService) GetChannelAccountById(channelAccountId uint) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	channelAccount, err := cas.channelAccountRepo.GetChannelAccountById(channelAccountId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	result["channel_account"] = channelAccount

	return result, nil
}

func (cas *ChannelAccountService) UpdateChannelAccount(ucam *presentation.UpdateChannelAccountModel) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	channelAccountId := ucam.ChannelAccountId
	newChannelAccount := entity.ChannelAccount{
		Name:                 ucam.Name,
		FaceboookPageId:      ucam.FaceboookPageId,
		InstagramId:          ucam.InstagramId,
		WhatsappNumId:        ucam.WhatsappNumId,
		WhatsappBusinessId:   ucam.WhatsappBusinessId,
		FacebookAccessToken:  ucam.FacebookAccessToken,
		InstagramAccessToken: ucam.InstagramAccessToken,
		WhatsappAccessToken:  ucam.WhatsappAccessToken,
	}

	channelAccount, err := cas.channelAccountRepo.UpdateChannelAccount(channelAccountId, &newChannelAccount)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	result["channel_account"] = channelAccount

	return result, nil
}

func (cas *ChannelAccountService) DeleteChannelAccountById(dcam *presentation.DeleteChannelAccountModel) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := cas.channelAccountRepo.DeleteChannelAccount(dcam.ChannelAccountId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	result["status"] = "SUCCESS"
	return result, nil
}
