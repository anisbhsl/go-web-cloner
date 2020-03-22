package scraper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

// Config contains the scraper configuration.
type Config struct {
	URL      string
	Includes []string
	Excludes []string

	ImageQuality uint // image quality from 0 to 100%, 0 to disable reencoding
	MaxDepth     uint // download depth, 0 for unlimited
	Timeout      uint // time limit in seconds to process each http request

	OutputDirectory     string
	Username            string
	Password            string
	ProjectID           string
	ScrapeID            string
	ScreenWidth         int      `json:"screen_width"`
	ScreenHeight        int      `json:"screen_height"`
	FolderThreshold     int      `json:"folder_threshold"`
	FolderExamplesCount int      `json:"folder_examples_count"`
	Patterns            []string `json:"patterns"`
}

// Scraper contains all scraping data.
type Scraper struct {
	config  Config
	log     *zap.Logger
	URL     *url.URL
	browser *browser.Browser

	cssURLRe *regexp.Regexp
	includes []*regexp.Regexp
	excludes []*regexp.Regexp

	// key is the URL of page or asset
	processed map[string]struct{}

	imagesQueue []*browser.DownloadableAsset
	report      *Report
}

// New creates a new Scraper instance.
func New(logger *zap.Logger, cfg Config) (*Scraper, error) {
	var errs *multierror.Error
	u, err := url.Parse(cfg.URL)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	includes, err := compileRegexps(cfg.Includes)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	excludes, err := compileRegexps(cfg.Excludes)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	if errs != nil {
		return nil, errs.ErrorOrNil()
	}

	if u.Scheme == "" {
		u.Scheme = "http" // if no URL scheme was given default to http
	}

	b := surf.NewBrowser()
	b.SetUserAgent(agent.GoogleBot())
	//b.SetTimeout(time.Duration(cfg.Timeout) * time.Second)

	s := &Scraper{
		config: cfg,

		browser:   b,
		log:       logger,
		processed: make(map[string]struct{}),
		URL:       u,
		cssURLRe:  regexp.MustCompile(`^url\(['"]?(.*?)['"]?\)$`),
		includes:  includes,
		excludes:  excludes,
		report: &Report{
			ProjectID:           cfg.ProjectID,
			ScrapeID:            cfg.ScrapeID,
			URL:                 cfg.URL,
			ScreenWidth:         cfg.ScreenWidth,
			ScreenHeight:        cfg.ScreenHeight,
			FolderThreshold:     cfg.FolderThreshold,
			FolderExamplesCount: cfg.FolderExamplesCount,
			Patterns:            cfg.Patterns,
			DetailedReport:      make([]DetailedReport, 0),
		},
	}

	return s, nil
}

// compileRegexps compiles the given regex strings to regular expressions
// to be used in the include and exclude filters.
func compileRegexps(sl []string) ([]*regexp.Regexp, error) {
	var errs error
	var l []*regexp.Regexp
	for _, e := range sl {
		re, err := regexp.Compile(e)
		if err == nil {
			l = append(l, re)
		} else {
			errs = multierror.Append(errs, err)
		}
	}
	return l, errs
}

// Start starts the scraping
func (s *Scraper) Start() error {
	if s.config.OutputDirectory != "" {
		if err := os.MkdirAll(s.config.OutputDirectory, os.ModePerm); err != nil {
			return err
		}
	}

	p := s.URL.Path
	if p == "" {
		p = "/"
	}
	s.processed[p] = struct{}{}

	if s.config.Username != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(s.config.Username + ":" + s.config.Password))
		s.browser.AddRequestHeader("Authorization", "Basic "+auth)
	}

	s.downloadPage(s.URL, 0, time.Now())

	err := s.generateReport()
	if err != nil {
		s.log.Error("report couldn't be generated",
			zap.Error(err))
	}

	return nil
}

func (s *Scraper) downloadPage(u *url.URL, currentDepth uint, startTime time.Time) {
	s.log.Info("Downloading", zap.Stringer("URL", u))
	if err := s.browser.Open(u.String()); err != nil {
		s.log.Error("Request failed",
			zap.Stringer("URL", u),
			zap.Error(err))

		currentTime := time.Now()
		details := DetailedReport{
			TimeStamp:    startTime,
			OriginURL:    u.Host + u.Path,
			LocalURL:     "",
			StatusCode:   400,
			ResponseTime: currentTime.Sub(startTime).Seconds(),
		}
		s.report.DetailedReport = append(s.report.DetailedReport, details)
		return
	}
	if c := s.browser.StatusCode(); c != http.StatusOK {
		s.log.Error("Request failed",
			zap.Stringer("URL", u),
			zap.Int("http_status_code", c))

		currentTime := time.Now()
		details := DetailedReport{
			TimeStamp:    startTime,
			OriginURL:    u.Host + u.Path,
			LocalURL:     "",
			StatusCode:   400,
			ResponseTime: currentTime.Sub(startTime).Seconds(),
		}
		s.report.DetailedReport = append(s.report.DetailedReport, details)

		return
	}

	buf := &bytes.Buffer{}
	if _, err := s.browser.Download(buf); err != nil {
		fmt.Println("[[151:scraper.go]] writing content to buffer hai guys!")
		s.log.Error("Downloading content failed",
			zap.Stringer("URL", u),
			zap.Error(err))

		currentTime := time.Now()
		details := DetailedReport{
			TimeStamp:    startTime,
			OriginURL:    u.Host + u.Path,
			LocalURL:     "",
			StatusCode:   400,
			ResponseTime: currentTime.Sub(startTime).Seconds(),
		}
		s.report.DetailedReport = append(s.report.DetailedReport, details)

		return
	}

	if currentDepth == 0 {
		u = s.browser.Url()
		// use the URL that the website returned as new base url for the
		// scrape, in case of a redirect it changed
		s.URL = u
	}

	s.storePage(u, buf, startTime)

	s.downloadReferences()

	var toScrape []*url.URL
	// check first and download afterwards to not hit max depth limit for
	// start page links because of recursive linking

	//TODO: this checks depth of URL hai guys
	for _, link := range s.browser.Links() {
		if s.checkPageURL(link.URL, currentDepth) {
			toScrape = append(toScrape, link.URL)
		}
	}

	//TODO: recursive downloading :)
	for _, URL := range toScrape {
		s.downloadPage(URL, currentDepth+1, time.Now())
	}
}

func (s *Scraper) storePage(u *url.URL, buf *bytes.Buffer, startTime time.Time) {
	var details DetailedReport
	html, err := s.fixFileReferences(u, buf)
	if err != nil {
		s.log.Error("Fixing file references failed",
			zap.Stringer("URL", u),
			zap.Error(err))

		currentTime := time.Now()
		details = DetailedReport{
			TimeStamp:    startTime,
			OriginURL:    u.Host + u.Path,
			LocalURL:     "",
			StatusCode:   400,
			ResponseTime: currentTime.Sub(startTime).Seconds(),
		}

	} else {
		buf = bytes.NewBufferString(html)
		filePath := s.GetFilePath(u, true)
		// always update html files, content might have changed
		if err = s.writeFile(filePath, buf); err != nil {
			s.log.Error("Writing HTML to file failed",
				zap.Stringer("URL", u),
				zap.String("file", filePath),
				zap.Error(err))
			currentTime := time.Now()
			details = DetailedReport{
				TimeStamp:    startTime,
				OriginURL:    u.Host + u.Path,
				LocalURL:     filePath,
				StatusCode:   400,
				ResponseTime: currentTime.Sub(startTime).Seconds(),
			}
		} else {
			currentTime := time.Now()
			details = DetailedReport{
				TimeStamp:    startTime,
				OriginURL:    u.Host + u.Path,
				LocalURL:     filePath,
				StatusCode:   200,
				ResponseTime: currentTime.Sub(startTime).Seconds(),
			}

		}
	}
	s.report.DetailedReport = append(s.report.DetailedReport, details)

}
