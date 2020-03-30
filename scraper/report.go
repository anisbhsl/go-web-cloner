package scraper

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Report struct {
	ProjectID           string           `json:"project_id"`
	ScrapeID            string           `json:"scrape_id"`
	URL                 string           `json:"url"`
	ScreenWidth         int              `json:"screen_width"`
	ScreenHeight        int              `json:"screen_height"`
	FolderThreshold     int              `json:"folder_threshold"`
	DetailedReport      []DetailedReport `json:"detailed_report"`
	FolderExamplesCount int              `json:"folder_examples_count"`
	Patterns            []string         `json:"patterns"`
}

type DetailedReport struct {
	TimeStamp    time.Time `json:"timestamp"`
	OriginURL    string    `json:"origin_url"`
	LocalURL     string    `json:"local_url"`
	StatusCode   int       `json:"status_code"`
	ResponseTime float64   `json:"response_time"`
}

func (s *Scraper) generateReport() error {
	file, err := json.Marshal(s.report)
	if err != nil {
		return err
	}

	//write report to a json file
	err = ioutil.WriteFile("data/"+s.report.ScrapeID+"_report.json", file, 0777)
	if err != nil {
		return err
	}
	s.log.Info("Report Generated...")
	return nil
}
