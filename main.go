package main

import (
	"github.com/NateshR/Likeminds-Real-Time/initialise"
)

func main() {

	initialise.AppConfig()
	initialise.ConnectDB()

	router := initialise.GetRouterEngine()
	initialise.StartServer(router)
}
