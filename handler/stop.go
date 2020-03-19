package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Stop(c *gin.Context){
	response:=make(map[string]interface{})
	response["msg"]="Scrapping Stopped"

	c.JSON(
		http.StatusOK,
		response,
	)

}