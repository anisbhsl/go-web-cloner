package asyncq

import (
	"go-web-cloner/scraper"
	"go.uber.org/zap"
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

func (d *Dispatcher) StartScrapper(s ScrapeConfig, scrapeID string) {
	cfg := scraper.Config{
		URL:                 s.URL,
		Includes:            nil,
		Excludes:            nil,
		ImageQuality:        0,
		MaxDepth:            0,
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
