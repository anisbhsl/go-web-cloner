package scraper

import (
	"encoding/json"
	"io/ioutil"
)

type Report struct{
	JobID           string `json:"job_id"`
	DetailedReport []DetailedReport  `json:"detailed_report"`
}

type DetailedReport struct{
	TimeStamp int64 `json:"timestamp"`
	OriginURL string `json:"origin_url"`
	LocalURL string `json:"local_url"`
	StatusCode int `json:"status_code"`
	ReponseTime int64 `json:"response_time"`
}

func (s *Scraper) generateReport() error {
	file, err := json.Marshal(s.report)
	if err!=nil{
		return err
	}
	//TODO: pass output directory plus path here
	err = ioutil.WriteFile("data/123random_report.json", file, 0777)
	if err!=nil{
		return err
	}
	return nil
}