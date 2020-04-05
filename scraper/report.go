package scraper

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
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
	tempReport:=make([]DetailedReport, 0)
	for _,val:=range s.report.DetailedReport{
		url:=val.OriginURL
		urlArr:=strings.Split(url,"/")
		length:=len(urlArr)
		path:=""
		if length<=2{  //its host only
			path=urlArr[0]+"/"
		}else{
			if urlArr[length-1]==""{
				for i:=0;i<length-2;i++{
					path+=urlArr[i]+"/"
				}
			}else{
				for i:=0;i<length-1;i++{
					path+=urlArr[i]+"/"
				}
			}
		}

		if v,ok:=s.Config.FolderCount[path];!ok{
			val.FolderCount="0"
		}else{
			count:=len(v)
			strCount:=strconv.Itoa(count)
			if count>=s.Config.FolderThreshold{
				strCount="{"+strCount+"}"
				val.FolderCount=strCount
			}else{
				val.FolderCount=strCount
			}
		}

		tempReport=append(tempReport,val)

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
	}

	s.report.DetailedReport=tempReport
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
