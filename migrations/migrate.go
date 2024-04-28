package main

import (
	"Omnichannel-CRM/domain/entity"
	"Omnichannel-CRM/package/config"
	"Omnichannel-CRM/package/database"
	"fmt"

	"Omnichannel-CRM/package/logger"

	"github.com/spf13/viper"
)

func init() {
	config.GetConfig()
}

func main() {
	// Init Log (should always be the first)
	logger.InitLogger()

	// Init database
	dbOmnichannel, err := database.InitDB(viper.GetString("Database.OmnichannelDBName"), viper.GetString("Database.Host"))
	if err != nil {
		logger.Error(fmt.Sprintf("[Migrations] Error when calling InitDB, trace: %+v", err))
		return
	}

	logger.Info("[InitDB] Migrating Database")
	err = dbOmnichannel.AutoMigrate(&entity.Interaction{})
	if err != nil {
		logger.Error(fmt.Sprintf("Error when migrating Interaction: trace: %+v", err))
		return
	}
	err = dbOmnichannel.AutoMigrate(&entity.Message{})
	if err != nil {
		logger.Error(fmt.Sprintf("Error when migrating Message: trace: %+v", err))
		return
	}

	err = dbOmnichannel.AutoMigrate(&entity.Reporter{})
	if err != nil {
		logger.Error(fmt.Sprintf("Error when migratingReporter: trace: %+v", err))
		return
	}

	err = dbOmnichannel.AutoMigrate(&entity.Thread{})
	if err != nil {
		logger.Error(fmt.Sprintf("Error when migrating Thread: trace: %+v", err))
		return
	}
}
