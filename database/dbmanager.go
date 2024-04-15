package database

import (
	"fmt"
	"log"
	"os"

	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	psql_con := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_NAME"))

	logFile, err := os.Create("gorm-log.txt")
	if err != nil {
		// Handle error
		log.Panic(err)
	}
	customLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags), // Use the log file as the output
		logger.Config{
			LogLevel:             logger.Info, // Set the log level (e.g., Info, Warn, Error, Silent)
			ParameterizedQueries: true,        // Don't include params in the SQL log
		},
	)
	db, err := gorm.Open(postgres.Open(psql_con), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		panic(err)
	}
	DB = db
	log.Println("Enabling UUID")
	createExtensionSQL := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
	if err := db.Exec(createExtensionSQL).Error; err != nil {
		log.Panic("failed to create uuid-ossp extension")
	}
	log.Println("Running Migrations")
	err = db.AutoMigrate(&structs.MotorPolicy{}, &structs.FileInformation{})
	if err != nil {
		log.Panic(err)
	}
}
