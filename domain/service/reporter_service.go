package service

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/presentation"

	"errors"

	"gorm.io/gorm"
)

type ReporterService struct {
	reporterRepo repository.IReporterRepository
}

type IReporterService interface {
	GetReporterByReporterId(uint) (map[string]interface{}, error)
	UpdateReporter(*presentation.UpdateReporterRequest) (map[string]interface{}, error)
}

func NewReporterService(reporterRepo repository.IReporterRepository) *ReporterService {
	reporterService := ReporterService{
		reporterRepo: reporterRepo,
	}
	return &reporterService
}

func (rs *ReporterService) GetReporterByReporterId(reporterId uint) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	reporter, err := rs.reporterRepo.GetReporterByReporterId(reporterId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	result["reporter"] = reporter

	return result, nil
}

func (rs *ReporterService) UpdateReporter(urr *presentation.UpdateReporterRequest) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	reporterId := urr.ReporterId
	newReporter := entity.Reporter{
		Name:             urr.Name,
		Email:            urr.Email,
		PhoneNumber:      urr.PhoneNumber,
		Gender:           urr.Gender,
		Address:          urr.Address,
		PlatformUsername: urr.PlatformUsername,
	}

	reporter, err := rs.reporterRepo.UpdateReporter(reporterId, &newReporter)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	result["reporter"] = reporter

	return result, nil
}
