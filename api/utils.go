package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrorQueryParamsMissing = "Query params missing!"
	ErrorBodyParamsMissing  = "Body params missing!"
)

type Response struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// GenerateResponse to be sent on request success
func GenerateResponse(c *gin.Context, dataResponse map[string]interface{}) {
	//Generating Response Object
	response := Response{
		Success: true,
	}

	//Removing Blank Data Key
	if len(dataResponse) > 0 {
		response.Data = dataResponse
	}

	c.JSON(http.StatusOK, response)
}

// GeneralAPIError to generate StatusInternalServerError
func GeneralAPIError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}

func GETQueryParamsMissingError(c *gin.Context) {
	GeneralBadRequestError(c, ErrorQueryParamsMissing)
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: ErrorQueryParamsMissing,
	})
}

func POSTBodyParamsMissingError(c *gin.Context) {
	GeneralBadRequestError(c, ErrorBodyParamsMissing)
}

func GeneralBadRequestError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}
