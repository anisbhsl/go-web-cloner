package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Status(c *gin.Context) {
	//Know if scrapper if idle or any process is running!
	response:=make(map[string]interface{})
	response["scrape_id"]="12345"
	response["msg"]="Scrapper Running"

	c.JSON(
		http.StatusOK,
		response,
	)
}
