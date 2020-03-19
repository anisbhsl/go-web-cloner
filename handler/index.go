package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {

		response:=make(map[string]interface{})
		response["msg"]="Website Cloner"
		response["status"]="Under Development : WIP"

		c.JSON(
			http.StatusOK,
			response,
		)
}
