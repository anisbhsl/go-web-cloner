package asyncq

import (
	"go-web-cloner/scraper"
	"go.uber.org/zap"
	"regexp"
	"strings"
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
	MaxDepth            uint     `json:"max_depth"`
}

func (d *Dispatcher) StartScrapper(s ScrapeConfig, scrapeID string) {
	cfg := scraper.Config{
		URL:                 s.URL,
		Includes:            nil,
		Excludes:            nil,
		ImageQuality:        0,
		MaxDepth:            s.MaxDepth,
		Timeout:             0,
		OutputDirectory:     "data/" + s.ProjectID + "/" + scrapeID, //set output directory
		Username:            s.Username,
		Password:            s.Password,
		ProjectID:           s.ProjectID,
		ScrapeID:            scrapeID,
		ScreenWidth:         s.ScreenWidth,
		ScreenHeight:        s.ScreenHeight,
		FolderThreshold:     s.FolderThreshold,
		FolderExamplesCount: s.FolderExamplesCount,
		Patterns:            s.Patterns,
		PatternCount:        initPatternCount(s.Patterns),
		FolderCount:         make(map[string]int),
		Stop:                false,
	}
	logger := logger()
	sc, err := scraper.New(logger, cfg)
	if err != nil {
		logger.Fatal("Initializing scraper failed", zap.Error(err))
	}

	logger.Info("Scraping", zap.Stringer("URL", sc.URL))
	d.Scraper = sc

	err = sc.Start()
	if err != nil {
		logger.Error("Scraping failed", zap.Error(err))
	}

	d.Queue = []string{} //dequeue

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

func initPatternCount(patterns []string)(map[*regexp.Regexp]int){
	patternCount:=make(map[*regexp.Regexp]int)

	if len(patterns)==0{
		//add main URL
		//pattern:=strings.Replace(url.Path,"*",".*",-1)
		//re:=regexp.MustCompile(pattern)
		//patternCount[re]=0
		return patternCount
	}else{
		for _,pa:=range patterns{
			if pa==""{
				continue
			}
			pattern:=strings.Replace(pa,"*",".*",-1)
			//remove http and https prefix
			if strings.Contains(pattern,"https://"){
				pattern=strings.TrimPrefix(pattern,"https://")
			}else if strings.Contains(pattern,"http://"){
				pattern=strings.TrimPrefix(pattern,"http://")
			}

			re:=regexp.MustCompile(pattern)
			patternCount[re]=0
		}
	}
	return patternCount
}