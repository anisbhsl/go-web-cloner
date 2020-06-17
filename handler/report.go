package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-web-cloner/scraper"
	"io/ioutil"
	"net/http"
	"os"
)

func Report(c *gin.Context) {
	scrapeID := c.Query("scrape_id")

	jsonFile, err := os.Open("data/" + scrapeID + "_report.json")
	if err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.tmpl",
			gin.H{
				"err": "no report found!",
			},
		)
		return
	}
	defer jsonFile.Close()

	reportInBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.tmpl",
			gin.H{
				"err": "Error while reading report",
			},
		)
		return
	}
	var finalReport scraper.Report
	_ = json.Unmarshal(reportInBytes, &finalReport)

	c.Writer.Header().Set("Access-Control-Allow-Origin","*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.HTML(
		http.StatusOK,
		"index.tmpl",
		gin.H{
			"scrape_id":   finalReport.ScrapeID,
			"project_id":finalReport.ProjectID,
			"screen_width":finalReport.ScreenWidth,
			"screen_height": finalReport.ScreenHeight,
			"folder_threshold":finalReport.FolderThreshold,
			"folder_examples_count":finalReport.FolderExamplesCount,
			"patterns":finalReport.Patterns,
			"details": finalReport.DetailedReport,
		},
	)

}
