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
	FolderCount   string      `json:"folder_count"`
}

func (s *Scraper) generateReport() error {
	/*
	Update folder counts for each pattern
	 */
	//tempReport:=make([]DetailedReport, 0)
	//for _,val:=range s.report.DetailedReport{
	//	if v:=s.Config.FolderCount[val.OriginURL];v==0{
	//		val.FolderCount="-"
	//	}else{
	//		strCount:=strconv.Itoa(s.Config.FolderThreshold)
	//		if s.Config.FolderCount[val.OriginURL]>=s.Config.FolderThreshold{
	//			strCount="{"+strCount+"}"
	//			val.FolderCount=strCount
	//		}else{
	//			val.FolderCount=strCount
	//		}
	//	}
	//
	//	tempReport=append(tempReport,val)
	//
	//}
	//
	//s.report.DetailedReport=tempReport
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
