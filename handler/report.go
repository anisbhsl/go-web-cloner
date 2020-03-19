package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Report(c *gin.Context) {

	response:=make(map[string]interface{})
	response["msg"]="Report under development"

	c.JSON(
		http.StatusOK,
		response,
	)

}
