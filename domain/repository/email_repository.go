package repository

import (
	"Omnichannel-CRM/package/utils"

	"google.golang.org/api/gmail/v1"
	"gorm.io/gorm"
)

type EmailRepository struct {
	db           *gorm.DB
	gmailService *gmail.Service
}

type IEmailRepository interface {
	GetHistoryList(historyId uint64) (res *gmail.ListHistoryResponse, err error)
	GetMessageById(messageId string) (res *gmail.Message, err error)
	SendEmail(message *gmail.Message) (res *gmail.Message, err error)
	GetProfile() (res *gmail.Profile, err error)
}

func NewEmailRepository(db *gorm.DB, gmailService *gmail.Service) *EmailRepository {
	emailRepo := EmailRepository{
		db:           db,
		gmailService: gmailService,
	}

	return &emailRepo
}

func (repo *EmailRepository) RefreshGmailService() {
	gmailService := utils.NewGmailService()
	repo.gmailService = gmailService
}

func (repo *EmailRepository) GetHistoryList(historyId uint64) (res *gmail.ListHistoryResponse, err error) {
	repo.RefreshGmailService()
	req := repo.gmailService.Users.History.List("me")
	req.StartHistoryId(historyId)
	res, err = req.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *EmailRepository) GetMessageById(messageId string) (res *gmail.Message, err error) {
	repo.RefreshGmailService()
	req := repo.gmailService.Users.Messages.Get("me", messageId)
	res, err = req.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *EmailRepository) SendEmail(message *gmail.Message) (res *gmail.Message, err error) {
	repo.RefreshGmailService()
	req := repo.gmailService.Users.Messages.Send("me", message)
	res, err = req.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *EmailRepository) GetProfile() (res *gmail.Profile, err error) {
	repo.RefreshGmailService()
	req := repo.gmailService.Users.GetProfile("me")
	res, err = req.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}
