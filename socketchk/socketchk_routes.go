package socketchk

import (
	"github.com/NateshR/Likeminds-Real-Time/middleware"
	"github.com/gin-gonic/gin"
)

// Router defines socket check paths
func Router(router *gin.Engine) {
	socketChkRoute := router.Group("/socket")
	{
		socketChkRoute.GET("", middleware.HttpConnectionUpgrader(), GetSocketStatus)
	}
}
