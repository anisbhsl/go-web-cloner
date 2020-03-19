package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Status(c *gin.Context) {
	response:=make(map[string]interface{})
	response["job_id"]="12345678"
	response["msg"]="Scrapper Running"

	c.JSON(
		http.StatusOK,
		response,
	)
}
