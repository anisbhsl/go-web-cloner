package handler

import (
	"github.com/gin-gonic/gin"
	"go-web-cloner/scraper"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type ScrapeConfig struct {
	URL                 string   `json:"url"`
	ScreenWidth         int      `json:"screen_width"`
	ScreenHeight        int      `json:"screen_height"`
	Username            string   `json:"username"`
	Password            string   `json:"password"`
	ProjectID           string   `json:"project_id,omitempty"`
	FolderThreshold     int      `json:"folder_threshold"`
	FolderExamplesCount int      `json:"folder_examples_count"`
	Patterns            []string `json:"patterns"`
}

//Scrape handles scrape request
func Scrape(c *gin.Context) {
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
	var scrapeConfig ScrapeConfig
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

	//send response
	curTime := time.Now().Unix()
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

	go startScrapper(scrapeConfig, scrapeID)

	c.JSON(
		http.StatusOK,
		response,
	)

}

func startScrapper(s ScrapeConfig, scrapeID string) {
	cfg := scraper.Config{
		URL:                 s.URL,
		Includes:            nil,
		Excludes:            nil,
		ImageQuality:        0,
		MaxDepth:            0,
		Timeout:             0,
		OutputDirectory:     "data/" + s.ProjectID+"/"+scrapeID, //set output directory
		Username:            s.Username,
		Password:            s.Password,
		ProjectID:           s.ProjectID,
		ScrapeID:            scrapeID,
		ScreenWidth:         s.ScreenWidth,
		ScreenHeight:        s.ScreenHeight,
		FolderThreshold:     s.FolderThreshold,
		FolderExamplesCount: s.FolderExamplesCount,
		Patterns:            s.Patterns,
	}
	logger := logger()
	sc, err := scraper.New(logger, cfg)
	if err != nil {
		logger.Fatal("Initializing scraper failed", zap.Error(err))
	}

	logger.Info("Scraping", zap.Stringer("URL", sc.URL))
	err = sc.Start()
	if err != nil {
		logger.Error("Scraping failed", zap.Error(err))
	}

}

//setup logger for web cloner
func logger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Development = false
	config.DisableCaller = true
	config.DisableStacktrace = true
	level := config.Level

	level.SetLevel(zap.InfoLevel)
	logger, _ := config.Build()
	return logger
}
