package scraper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"go-web-cloner/surf"
	"go-web-cloner/surf/agent"
	"go-web-cloner/surf/browser"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Config contains the scraper configuration.
type Config struct {
	URL      string
	Includes []string
	Excludes []string

	ImageQuality uint
	MaxDepth     uint // download depth, 0 for unlimited
	Timeout      uint // time limit in seconds to process each http request

	OutputDirectory     string
	Username            string
	Password            string
	AccessToken         string
	ProjectID           string
	ScrapeID            string
	ScreenWidth         int      `json:"screen_width"`
	ScreenHeight        int      `json:"screen_height"`
	FolderThreshold     int      `json:"folder_threshold"`
	FolderExamplesCount int      `json:"folder_examples_count"`
	Patterns            []string `json:"patterns"`
	PatternCount map[*regexp.Regexp]int //folder count patterns
	FolderCount  map[string]map[string]bool //Folder Threshold Count Global
	MandatoryDownload map[string]bool //download links which have been resolved due to folder count limit
	Stop bool
}

// Scraper contains all scraping data.
type Scraper struct {
	Config  Config
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
		Config: cfg,

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
			//ScreenWidth:         cfg.ScreenWidth,
			//ScreenHeight:        cfg.ScreenHeight,
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
	//create output DIR
	if s.Config.OutputDirectory != "" {
		if err := os.MkdirAll(s.Config.OutputDirectory, os.ModePerm); err != nil {
			return err
		}
	}

	p := s.URL.Path
	if p == "" {
		p = "/"
	}
	s.processed[p] = struct{}{}

	if s.Config.Username != "" && s.Config.Password!=""{
		s.log.Info("Setting request header for",zap.String("basic auth",s.Config.Username))

		auth := base64.StdEncoding.EncodeToString([]byte(s.Config.Username + ":" + s.Config.Password))
		s.browser.AddRequestHeader("Authorization", "Basic "+auth)
		s.log.Info("basic auth encoded: ",zap.String("auth: ",auth))
	}

	//if s.Config.AccessToken!=""{
	//	s.log.Info("Setting request header ",zap.String("access token",s.Config.AccessToken))
	//	s.browser.AddRequestHeader("Authorization","Bearer "+s.Config.AccessToken)
	//}

	s.downloadPage(s.URL, 0, time.Now())

	err := s.generateReport()
	if err != nil {
		s.log.Error("report couldn't be generated",
			zap.Error(err))
	}

	return nil
}

func (s *Scraper) downloadPage(u *url.URL, currentDepth uint, startTime time.Time) {
	if s.Config.Stop{
		s.log.Info("Stop Command Received ==> Stopping scrapper...")
		return
	}
	/*
		Check folder count here:
		if a/b/c/x/index.html --> get path a/b/c/x only and break it into a/b/c

		if a/b/c exists,
			map["a/b/c"]++
		else
		  map["a/b/c"]=1
	*/
	log.Println("Downloading URL path: ",u.Scheme+"://"+u.Host+u.Path)
	println("contents in mandatory download")
	for k,_:=range s.Config.MandatoryDownload{
		println(k)
	}
	if _,ok:=s.Config.MandatoryDownload[u.Scheme+"://"+u.Host+u.Path];!ok{
		println("doesn't exist in mandatory download so check for folder count")
		//if folder count threshold has exceeded
		//do not visit any links just return
		if s.hasFolderCountExceeded(u){
			return
		}

		//if folder threshold count has been exceeded return
		//do not visit any links inside
		if s.hasFolderThresholdExceededForPattern(u){
			return
		}
	}else{
		delete(s.Config.MandatoryDownload,u.Scheme+"://"+u.Host+u.Path) //delete matched urls
	}


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

	if c := s.browser.StatusCode(); c !=http.StatusOK {

		s.log.Info("connection header: ",zap.String("header=>",s.browser.ResponseHeaders().Get("Connection")))

		s.log.Error("Request failed",
			zap.Stringer("URL", u),
			zap.Int("http_status_code", c))

		currentTime := time.Now()
		details := DetailedReport{
			TimeStamp:    startTime,
			OriginURL:    u.Host + u.Path,
			LocalURL:     "",
			StatusCode:   c,
			ResponseTime: currentTime.Sub(startTime).Seconds(),
		}
		s.report.DetailedReport = append(s.report.DetailedReport, details)

		return
	}


	buf := &bytes.Buffer{}
	if _, err := s.browser.Download(buf); err != nil {
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

	for _, link := range s.browser.Links() {
		if s.checkPageURL(link.URL, currentDepth) {
			toScrape = append(toScrape, link.URL)
		}
	}

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
		f := fmt.Sprintf("filepath is : %v", filePath)
		upath := fmt.Sprintf("u.Path is : %v", u.Path)
		s.log.Info(f)
		s.log.Info(upath)

		regexForResources := regexp.MustCompile("(\\.png|\\.jpg|\\.jpeg|\\.pdf|\\.gif|\\.docx|\\.mp4|\\.avi)")
		resourceExtension := regexForResources.FindString(u.Path)

		if resourceExtension != "" {  //if resource extension is other than .html do not download it's equivalent html file
			return
		}

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
				LocalURL:     "", //TODO: add resource mapping logic here
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
		s.report.DetailedReport = append(s.report.DetailedReport, details)

	}
}

//checks folder count for specified patterns only
func (s *Scraper) hasFolderThresholdExceededForPattern(u *url.URL) bool{

	/*
	TODO:
	1. Check if the URL is in patterns
	2. If yes, increase its coutner and check if counter has exceeded.
		2.i. If YES, add the URL pattern to exclude list
		s.Excludes=append(s.Excludes,url)
	*/

	//if there are no patterns supplied, do job as usual
	if len(s.Config.Patterns)==0{
		return false
	}

	for reg,_:=range s.Config.PatternCount{
		s.log.Info("[[hasFolderThresholdExceeded]] u.path+u.host:",zap.String("url: %v",u.Host+u.Path))
		if reg.Match([]byte(u.Host+u.Path)){
		    s.log.Info("matched with",zap.String("url : ",u.Host+u.Path))
		    s.Config.PatternCount[reg]++ //increase count for matched
			if s.Config.FolderThreshold < s.Config.PatternCount[reg]  {
				return true
			}
			return false
		}
	}

	return false
}

func (s *Scraper) hasFolderCountExceeded(u *url.URL) bool{   //checks folder count globally
	if s.Config.FolderThreshold==0{
		return false
	}

	host:=u.Host
	path:=u.Path

	if path=="" || path=="/"{   //no need to check initially when path is empty
		return false
	}
	s.log.Info("path",zap.String("path: ",path))
	pathArr:=strings.Split(path,"/")

	length:=len(pathArr)
	s.log.Info("PathArr --->",zap.Strings("pathArr:",pathArr))

	if length<2 {                 //TODO: work on this logic here man
		return false
	}
	s.log.Info(" PathArr ----> ",zap.Strings("pathArr:",pathArr))

	newPathArr:=[]string{} //will hold path array only
	for _,val:=range pathArr{
		if val!=""{
			newPathArr=append(newPathArr,val)
		}
	}
	s.log.Info("newPathArr: ",zap.Strings("newPathArr:",newPathArr))
	//if pathArr[length-1]!=""{  //if there is no slash at last, the html page is in same dir level


	finalPath:=host+"/"
	//for host part
	if val,ok:=s.Config.FolderCount[finalPath];!ok{
		folder:=make(map[string]bool)
		folder[newPathArr[0]]=true

		s.Config.FolderCount[finalPath]=folder
	}else{
		l:=len(val)
		if l>=s.Config.FolderThreshold{
			if _,o:=val[newPathArr[0]];!o{
				s.log.Info("folder threshold reached")
				return true
			}
		}else{ //todo: yo else chainna hola yar
			//if not add folder
			//check for / at the end of path
			if len(newPathArr)==1{
				if pathArr[len(pathArr)-1]==""{
					val[newPathArr[0]+"/"]=true
					s.Config.FolderCount[finalPath]=val//update folders
				}else{
					val[newPathArr[0]]=true
					s.Config.FolderCount[finalPath]=val //update folders
				}
			}else{
				val[newPathArr[0]]=true
				s.Config.FolderCount[finalPath]=val //update folders
			}
		}
	}

	//for later path
	length=len(newPathArr)
	for i:=0;i<length-1;i++{

		finalPath+=newPathArr[i]
		finalPath+="/"

		if val,ok:=s.Config.FolderCount[finalPath];!ok{
			folder:=make(map[string]bool)

			folder[newPathArr[i+1]]=true
			s.Config.FolderCount[finalPath]=folder

		}else{
			l:=len(val)

			if l>=s.Config.FolderThreshold{
				if _,o:=val[newPathArr[i+1]];o{
					continue
				}
				s.log.Info("folder threshold reached")
				return true
			}

			//if not add folder
			if i==length-2{
				if pathArr[len(pathArr)-1]==""{
					val[newPathArr[i+1]+"/"]=true
					s.Config.FolderCount[finalPath]=val//update folders
				}else{
					val[newPathArr[i+1]]=true
					s.Config.FolderCount[finalPath]=val //update folders
				}

			}else{
				val[newPathArr[i+1]]=true
				s.Config.FolderCount[finalPath]=val //update folders
			}
		}
	}
	return false
}