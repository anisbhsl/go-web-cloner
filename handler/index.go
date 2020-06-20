package handler

import (
	"github.com/gin-gonic/gin"
	asyncq "go-web-cloner/asynq"
	"net/http"
)

func Index(dispatcher *asyncq.Dispatcher) gin.HandlerFunc{
	return func(c *gin.Context){
		response := make(map[string]interface{})
		response["msg"] = "Website Cloner"
		response["status"] = "Under Development : WIP"
		c.JSON(
			http.StatusOK,
			response,
		)
	}
}
