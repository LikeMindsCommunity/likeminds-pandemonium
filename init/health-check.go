package init

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck | HealthCheck is used to validate health check and host / path
func HealthCheck(c *gin.Context) {
	//Send response with success as true
	c.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}
