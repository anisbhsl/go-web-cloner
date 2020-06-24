package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-cloner/scraper"
	"io/ioutil"
	"net/http"
	"os"
)

func Report(c *gin.Context) {
	scrapeID := c.Query("scrape_id")
	response:=make(map[string]interface{})
	jsonFile, err := os.Open("data/" + scrapeID + "_report.json")
	if err != nil {
		//c.HTML(
		//	http.StatusBadRequest,
		//	"error.tmpl",
		//	gin.H{
		//		"err": "no report found!",
		//	},
		//)
		response["err"]=fmt.Sprintf("no report found")
		c.JSON(
			http.StatusBadRequest,
			response,
			)
		return
	}
	defer jsonFile.Close()

	reportInBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		//c.HTML(
		//	http.StatusBadRequest,
		//	"error.tmpl",
		//	gin.H{
		//		"err": "Error while reading report",
		//	},
		//)
		response["err"]=fmt.Sprintf("error while reading report")
		c.JSON(
			http.StatusBadRequest,
			response,
			)

		return
	}
	var finalReport scraper.Report
	_ = json.Unmarshal(reportInBytes, &finalReport)

	response["report"]=finalReport

	//c.HTML(
	//	http.StatusOK,
	//	"index.tmpl",
	//	gin.H{
	//		"scrape_id":   finalReport.ScrapeID,
	//		"project_id":finalReport.ProjectID,
	//		"screen_width":finalReport.ScreenWidth,
	//		"screen_height": finalReport.ScreenHeight,
	//		"folder_threshold":finalReport.FolderThreshold,
	//		"folder_examples_count":finalReport.FolderExamplesCount,
	//		"patterns":finalReport.Patterns,
	//		"details": finalReport.DetailedReport,
	//	},
	//)
	c.JSON(
		http.StatusOK,
		response,
		)
}
