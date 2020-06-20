package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	asyncq "go-web-cloner/asynq"
	"net/http"
)


func Stop(dispatcher *asyncq.Dispatcher) gin.HandlerFunc{
	return func(c *gin.Context){
		response:=make(map[string]interface{})

		if ok:=dispatcher.IsWorkerAvailable();ok{
			response["msg"]="No active jobs"
			c.JSON(
				http.StatusOK,
				response,
				)
			return
		}

		dispatcher.StopScrapper() //stop running job

		response["msg"]=fmt.Sprintf("Scrapping Stopped for scrape_id: %v",dispatcher.Queue)

		c.JSON(
			http.StatusOK,
			response,
		)
	}
}
