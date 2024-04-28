package repository

import (
	"Omnichannel-CRM/domain/entity"

	"gorm.io/gorm"
)

type ReporterRepository struct {
	db *gorm.DB
}

type IReporterRepository interface {
	CreateReporter(*entity.Reporter) (*entity.Reporter, error)
	GetReporterByMetaReporterId(string) (*entity.Reporter, error)
	GetReporterByReporterId(uint) (*entity.Reporter, error)
	GetReporterByEmail(string) (*entity.Reporter, error)
	UpdateReporter(uint, *entity.Reporter) (*entity.Reporter, error)
}

func NewReporterRepository(db *gorm.DB) *ReporterRepository {
	reporterRepo := ReporterRepository{
		db: db,
	}
	return &reporterRepo
}

func (rr *ReporterRepository) CreateReporter(reporter *entity.Reporter) (*entity.Reporter, error) {
	err := rr.db.Create(&reporter).Error

	if err != nil {
		return nil, err
	}

	return reporter, nil
}

func (rr *ReporterRepository) GetReporterByMetaReporterId(metaReporterId string) (*entity.Reporter, error) {
	var reporter entity.Reporter

	err := rr.db.Where("meta_reporter_id = ?", metaReporterId).Take(&reporter).Error

	if err != nil {
		return nil, err
	}
	return &reporter, nil
}

func (rr *ReporterRepository) GetReporterByReporterId(reporterId uint) (*entity.Reporter, error) {
	var reporter entity.Reporter

	err := rr.db.Where("id = ?", reporterId).Take(&reporter).Error

	if err != nil {
		return nil, err
	}
	return &reporter, nil
}

func (rr *ReporterRepository) GetReporterByEmail(email string) (*entity.Reporter, error) {
	var reporter entity.Reporter

	err := rr.db.Where("email = ?", email).Take(&reporter).Error

	if err != nil {
		return nil, err
	}
	return &reporter, nil
}

func (rr *ReporterRepository) UpdateReporter(reporterId uint, newReporter *entity.Reporter) (*entity.Reporter, error) {
	var currentReporter entity.Reporter

	err := rr.db.Where("id", reporterId).Take(&currentReporter).Error

	if err != nil {
		return nil, err
	}

	if newReporter.Name != "" {
		currentReporter.Name = newReporter.Name
	}
	if newReporter.Email != "" {
		currentReporter.Email = newReporter.Email
	}
	if newReporter.PhoneNumber != "" {
		currentReporter.PhoneNumber = newReporter.PhoneNumber
	}
	if newReporter.Gender != "" {
		currentReporter.Gender = newReporter.Gender
	}
	if newReporter.Address != "" {
		currentReporter.Address = newReporter.Address
	}
	if newReporter.PlatformUsername != "" {
		currentReporter.PlatformUsername = newReporter.PlatformUsername
	}

	err = rr.db.Save(&currentReporter).Error
	if err != nil {
		return nil, err
	}

	return &currentReporter, nil
}
