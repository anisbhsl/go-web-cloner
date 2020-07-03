package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	asyncq "go-web-cloner/asynq"
	"go-web-cloner/scraper"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//Report generates report for given scrapeID
func Report(dispatcher *asyncq.Dispatcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		scrapeID := c.Query("scrape_id")
		format := c.Query("format")
		//if format is not specified return HTML response
		if strings.ToLower(format) != "json" {
			generateHTMLReport(c, scrapeID)
			return
		}

		response := make(map[string]interface{})
		if !dispatcher.IsWorkerAvailable(){
			if dispatcher.Scraper.Config.ScrapeID==scrapeID{
				response["report"]=dispatcher.Scraper.GenerateDynamicReport()
				c.JSON(
					http.StatusPartialContent,
					response,
					)
				return
			}
		}

		jsonFile, err := os.Open("data/" + scrapeID + "_report.json")
		if err != nil {
			response["err"] = fmt.Sprintf("no report found")
			c.JSON(
				http.StatusBadRequest,
				response,
			)
			return
		}
		defer jsonFile.Close()

		reportInBytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			response["err"] = fmt.Sprintf("error while reading report")
			c.JSON(
				http.StatusBadRequest,
				response,
			)

			return
		}
		var finalReport scraper.Report
		_ = json.Unmarshal(reportInBytes, &finalReport)

		response["report"] = finalReport
		c.JSON(
			http.StatusOK,
			response,
		)
	}
}

//generateHTMLReport generates a HTML Report for the given scrapeID
func generateHTMLReport(c *gin.Context,scrapeID string){
	jsonFile, err := os.Open("data/" + scrapeID + "_report.json")
	if err!=nil{
		c.HTML(
			http.StatusBadRequest,
			"error.tmpl",
			gin.H{
				"err": "no report found!",
			},
		)
	}

	defer jsonFile.Close()
	reportInBytes, err := ioutil.ReadAll(jsonFile)
	if err!=nil{
		c.HTML(
			http.StatusBadRequest,
			"error.tmpl",
			gin.H{
				"err": "Error while reading report",
			},
		)
	}

	var finalReport scraper.Report
	_ = json.Unmarshal(reportInBytes, &finalReport)

	c.HTML(
		http.StatusOK,
		"index.tmpl",
		gin.H{
			"scrape_id":   finalReport.ScrapeID,
			"project_id":finalReport.ProjectID,
			"folder_threshold":finalReport.FolderThreshold,
			"folder_examples_count":finalReport.FolderExamplesCount,
			"patterns":finalReport.Patterns,
			"details": finalReport.DetailedReport,
		},
	)

}