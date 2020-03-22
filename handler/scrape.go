package handler

import (
	"github.com/gin-gonic/gin"
	asyncq "go-web-cloner/asynq"
	"net/http"
	"strconv"
	"time"
)

//Scrape handles scrape request
func Scrape(dispatcher *asyncq.Dispatcher) gin.HandlerFunc{
	return func(c *gin.Context) {
		/*
		   Sample POST request body
		   {
		       "url":"www.airbnb.com/hosting",
		       "screen_width":1920,
		       "screen_height": 1080,
		       "username": "abc@def.gh",
		       "password": "L36gh!h'",
		        "project_id": "abc", //optional
		       "folder_threshold": 20,
		        "folder_examples_count":3,
		       "patterns": ["www.airbnb.com/s/asterisk(*)/experiences]"

		*/
		var scrapeConfig asyncq.ScrapeConfig
		if err := c.Bind(&scrapeConfig); err != nil {
			response := make(map[string]interface{})
			response["err"] = "Error while form parsing"
			response["err_desc"] = err
			c.JSON(
				http.StatusBadRequest,
				response,
			)
			return
		}

		if ok:=dispatcher.IsWorkerAvailable();!ok{
			response := make(map[string]interface{})
			response["msg"]="Scrapper Running Another Job"
			response["scrape_id"]=dispatcher.Queue
			c.JSON(
				http.StatusTooManyRequests,
				response,
				)
			return
		}

		//send response
		curTime:= time.Now().Unix()
		scrapeID := strconv.Itoa(int(curTime))
		response := make(map[string]interface{})
		response["scrape_id"] = scrapeID //use unix timestamp
		response["msg"] = "Scrapping Started"
		response["url"] = scrapeConfig.URL
		response["screen_height"] = scrapeConfig.ScreenHeight
		response["screen_width"] = scrapeConfig.ScreenWidth
		response["project_id"] = scrapeConfig.ProjectID
		if response["project_id"] == "" {
			response["project_id"] = "default"
			scrapeConfig.ProjectID = "default"
		}
		response["folder_threshold"] = scrapeConfig.FolderExamplesCount
		response["folder_examples_count"] = scrapeConfig.FolderExamplesCount
		response["patterns"] = scrapeConfig.Patterns

		//go startScrapper(scrapeConfig, scrapeID)
		dispatcher.Queue=append(dispatcher.Queue,scrapeID) //enqueue a job
		go dispatcher.StartScrapper(scrapeConfig,scrapeID)  //run async

		c.JSON(
			http.StatusOK,
			response,
		)

	}
}

