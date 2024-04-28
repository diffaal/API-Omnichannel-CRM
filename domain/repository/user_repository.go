package repository

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/package/enum"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	GetUserListByIds([]string) ([]entity.User, error)
	GetAgentList(map[string]interface{}) ([]entity.User, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	userRepo := UserRepository{
		db: db,
	}
	return &userRepo
}

func (ur *UserRepository) GetUserListByIds(userIds []string) ([]entity.User, error) {
	var userList []entity.User

	result := ur.db.Where("id IN ?", userIds).Find(&userList)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return userList, nil
}

func (ur *UserRepository) GetAgentList(filters map[string]interface{}) ([]entity.User, error) {
	var agentList []entity.User
	roles := []int{enum.ROLE_ADMIN_PUSAT, enum.ROLE_AGENT_PUSAT}

	queryDB := ur.db.Where("role in ?", roles)

	if filters["agent_ids"] != nil {
		queryDB = queryDB.Where("id IN ?", filters["agent_ids"])
	}
	if filters["keyword"] != nil {
		queryDB = queryDB.Where("first_name LIKE ? OR last_name LIKE ?", fmt.Sprintf("%%%s%%", filters["keyword"]))
	}
	if filters["page"] != nil && filters["pageSize"] != nil {
		offset := (filters["page"].(int) - 1) * filters["pageSize"].(int)
		limit := filters["pageSize"].(int)
		queryDB = queryDB.Offset(offset).Limit(limit)
	}

	result := queryDB.Find(&agentList)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return agentList, nil
}
