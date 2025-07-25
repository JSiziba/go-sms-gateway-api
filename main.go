package main

import (
	"fmt"
	"go-sms-gateway-api/config"
	_ "go-sms-gateway-api/docs"
	"go-sms-gateway-api/models"
	"go-sms-gateway-api/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

// @title SMS Gateway
// @version 1.0
// @description SMS Gateway Server
// @host localhost:1949
// @BasePath /api/v1
// @securityDefinitions.apikey X-Require-Whisk-Auth
// @in header
// @name X-Require-Whisk-Auth
// @description Enter your API key
func main() {
	apiArtClean := `
        AAAAA         PPPPPP       IIIIIII
       A     A        P    PP        II
      A       A       P     PP       II
     AAAAAAAAA        PPPPPP         II
    A         A       P              II
   A           A      P            IIIIIII
  -------------------------------------------
      S E R V E R   S T A R T I N G . . .
  -------------------------------------------`
	fmt.Println(apiArtClean)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.SMSMessage{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	srv := server.NewServer(db, &cfg)

	port := strconv.Itoa(cfg.ServerPort)
	log.Printf("\033[32mServer starting on port %s\033[0m", port)
	err = srv.Start(":" + port)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
