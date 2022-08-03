package initialise

import (
	"log"

	"github.com/NateshR/Likeminds-Real-Time/utility"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StartServer starts a http webserver and attaches routergroup
func StartServer(router *gin.Engine) {

	appPort, ok := viper.Get(EnvVarAppPort).(string)
	if !ok {
		log.Println("Application run port missing. Starting on default port..")
	}

	if len(appPort) == 0 {
		appPort = AppDefaultPort
	}

	portStrList := []string{":", appPort}
	joinChar := ""
	runStr := utility.JoinStrs(portStrList, joinChar)

	log.Fatal(router.Run(runStr))
}
