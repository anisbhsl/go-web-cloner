package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	asyncq "go-web-cloner/asynq"
	"net/http"
)

func Status(dispatcher *asyncq.Dispatcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := make(map[string]interface{})

		//Know if scrapper if idle or any process is running!
		if ok:=dispatcher.IsWorkerAvailable();!ok{
			response["status"]=fmt.Sprintf("Another scrapper with id %v already running",dispatcher.Queue)
			c.JSON(
				http.StatusTooManyRequests,
				response,
				)
			return
		}

		response["status"] = "IDLE"
		c.JSON(
			http.StatusOK,
			response,
		)
	}
}
