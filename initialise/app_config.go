package initialise

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// AppConfig sets application config
func AppConfig() {
	setAppEnvVars()
	setAppMode()
}

func setAppEnvVars() {

	viper.SetConfigFile(EnvConfigFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading config file, reason=%s \nExiting..", err))
	}
}

func setAppMode() {

	appMode, ok := viper.Get(EnvVarAppMode).(string)
	if !ok {
		log.Fatal("Application run mode missing. Exiting..")
	}

	switch appMode {
	case AppModeRelease:
		appMode = gin.ReleaseMode
	case AppModeTest:
		appMode = gin.TestMode
	case AppModeDebug:
		appMode = gin.DebugMode
	default:
		log.Fatal(fmt.Sprintf("Application run mode unknown, run_mode=%s \n(available mode: debug release test)", appMode))
	}

	gin.SetMode(appMode)
}
