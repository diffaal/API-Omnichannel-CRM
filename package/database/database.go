package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(DBName string, host string) (*gorm.DB, error) {
	username := viper.GetString("Database.Username")
	password := viper.GetString("Database.Password")
	port := viper.GetInt("Database.Port")
	dbname := DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, username, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
