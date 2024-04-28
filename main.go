package main

import (
	"Omnichannel-CRM/api"
	"Omnichannel-CRM/package/config"
	"Omnichannel-CRM/package/database"
	"Omnichannel-CRM/package/logger"
	"flag"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func init() {
	config.GetConfig()
}

func main() {
	logger.InitLogger()

	dbCRM, err := database.InitDB(viper.GetString("Database.CRMDBName"), viper.GetString("Database.CRMDBHost"))
	if err != nil {
		log.Fatal(err)
		return
	}

	dbOmnichannel, err := database.InitDB(viper.GetString("Database.OmnichannelDBName"), viper.GetString("Database.Host"))
	if err != nil {
		log.Fatal(err)
		return
	}

	var port int
	flag.IntVar(&port, "port", viper.GetInt("App.Port"), "Port to run the server on")
	flag.Parse()

	app := api.SetupRouter(dbCRM, dbOmnichannel)
	app.Run(fmt.Sprintf(":%d", port))
}
