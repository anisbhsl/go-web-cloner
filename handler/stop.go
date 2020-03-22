package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Stop(c *gin.Context){
	response:=make(map[string]interface{})
	response["msg"]="Scrapping Stopped"
	//generate report upto this !
	c.JSON(
		http.StatusOK,
		response,
	)

}