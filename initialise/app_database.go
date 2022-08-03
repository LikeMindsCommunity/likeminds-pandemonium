package initialise

import (
	"database/sql"
	"fmt"
	"log"

	// blank import used to import the postgres library
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	db *sql.DB
)

// ConnectDB return a database connection
func ConnectDB() {

	var (
		host     = getDbVarFromEnv(EnvVarDbHost)
		port     = getDbVarFromEnv(EnvVarDbPort)
		dbname   = getDbVarFromEnv(EnvVarDbName)
		user     = getDbVarFromEnv(EnvVarDbUser)
		password = getDbVarFromEnv(EnvVarDbPwd)
	)

	psqlConfig := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlConfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot open database connection. Failed, reason=%s", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot connect with database. Failed, reason=%s", err))
	}
}

func getDbVarFromEnv(EnvVarDbVar string) string {
	dbVar, ok := viper.Get(EnvVarDbVar).(string)
	if !ok {
		log.Fatal(fmt.Sprintf("Database parameter %s missing. Exiting..", EnvVarDbVar))
	}

	return dbVar
}
