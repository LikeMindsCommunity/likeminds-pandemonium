package statuschk

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAppStatus returns application status
func GetAppStatus(c *gin.Context) {
	const msgStr string = "Application server up and running.."
	const ErrStr string = ""

	var apiresponse Response = Response{
		Status: http.StatusOK,
		Data:   msgStr,
		Error:  ErrStr,
	}

	c.JSON(apiresponse.Status, apiresponse)
}
