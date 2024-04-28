package entity

import (
	"Omnichannel-CRM/package/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"unique" json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `gorm:"unique" json:"email"`
	Role      int    `json:"role"`
	Province  string `json:"province"`
	City      string `json:"city"`
}

func (u *User) EncryptPassword(password string) (string, error) {
	hashPassword, err := security.Hash(password)
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil

}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return nil
}
