package repository

import (
	"Omnichannel-CRM/domain/entity"

	"gorm.io/gorm"
)

type ThreadRepository struct {
	db *gorm.DB
}

type IThreadRepository interface {
	GetThreadByID(threadId string) (thread entity.Thread, err error)
	InsertThread(thread entity.Thread) (entity.Thread, error)
}

func NewThreadRepository(db *gorm.DB) *ThreadRepository {
	threadRepo := ThreadRepository{
		db: db,
	}

	return &threadRepo
}

func (repo *ThreadRepository) GetThreadByID(threadId string) (thread entity.Thread, err error) {
	err = repo.db.First(&thread, "id = ?", threadId).Error
	if err != nil {
		return thread, err
	}

	return thread, nil
}

func (repo *ThreadRepository) InsertThread(thread entity.Thread) (entity.Thread, error) {
	err := repo.db.Create(&thread).Error
	if err != nil {
		return thread, err
	}

	return thread, nil
}
