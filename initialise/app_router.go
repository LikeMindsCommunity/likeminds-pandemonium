package initialise

import (
	"github.com/NateshR/Likeminds-Real-Time/socketchk"
	"github.com/NateshR/Likeminds-Real-Time/statuschk"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AppRouter defines router struct
type AppRouter struct {
	router *gin.Engine
}

// GetRouterEngine returns gin engine router
func GetRouterEngine() *gin.Engine {

	appRouter := AppRouter{
		router: gin.Default(),
	}

	addTrustedProxies(appRouter.router)
	addRouterCors(appRouter.router)
	addRouterRoutes(appRouter)

	return appRouter.router

}

func addTrustedProxies(router *gin.Engine) {
	router.SetTrustedProxies(nil)
}

func addRouterCors(router *gin.Engine) {
	router.Use(cors.New(enableCors()))
}

func enableCors() cors.Config {

	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true

	conf.AddAllowHeaders(
		AppAllowedHeaderAPIKey,
		AppAllowedHeaderAuthorisation,
		AppAllowedHeaderDeviceID,
		AppAllowedHeaderMemberID,
		AppAllowedHeaderPassword,
		AppAllowedHeaderPlatformCode,
		AppAllowedHeaderPlatformType,
		AppAllowedHeaderUserName,
		AppAllowedHeaderVersionCode,
	)

	return conf
}

func addRouterRoutes(appRouter AppRouter) {
	statuschk.Router(appRouter.router)
	socketchk.Router(appRouter.router)
}
