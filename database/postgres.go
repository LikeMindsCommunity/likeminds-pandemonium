package database

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/common"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() *gorm.DB {

	host := common.GoDotEnvVariable(common.DotEnvVarPostgresHost)
	port := common.GoDotEnvVariable(common.DotEnvVarPostgresPort)
	user := common.GoDotEnvVariable(common.DotEnvVarPostgresUser)
	password := common.GoDotEnvVariable(common.DotEnvVarPostgresUserPassword)
	database := common.GoDotEnvVariable(common.DotEnvVarPostgresDatabaseName)
	sslMode := "allow"
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
