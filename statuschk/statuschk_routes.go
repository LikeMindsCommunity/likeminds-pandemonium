package statuschk

import "github.com/gin-gonic/gin"

// Router defines status check paths
func Router(router *gin.Engine) {
	statusChkRoute := router.Group("/status")
	{
		statusChkRoute.GET("", GetAppStatus)
	}
}
