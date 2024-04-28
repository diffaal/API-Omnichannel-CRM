package entity

import "gorm.io/gorm"

type Thread struct {
	gorm.Model
	ID        string `json:"id" gorm:"primaryKey"`
	Subject   string `json:"subject"`
	EmailDate string `json:"email_date"`
	From      string `json:"from"`
}
