package database

import (
	"fmt"
	"likeminds-pandemonium/common"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() *gorm.DB {

	host := common.GoDotEnvVariable("POSTGRES_HOST")
	port := common.GoDotEnvVariable("POSTGRES_PORT")
	user := common.GoDotEnvVariable("POSTGRES_USER")
	password := common.GoDotEnvVariable("POSTGRES_USER_PASSWORD")
	database := common.GoDotEnvVariable("POSTGRES_DATABASE_NAME")
	sslMode := "disable"
	timeZone := "Asia/Kolkata"

	dsn := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s  sslmode=%s TimeZone=%s",
		host,
		port,
		user,
		password,
		database,
		sslMode,
		timeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database = %s", database)
	}

	log.Printf("successfully connected to database = %s", database)
	return db
}

var (
	Postgres *gorm.DB
)
