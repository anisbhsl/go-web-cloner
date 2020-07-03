package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"io"
	"net/url"
	"strings"
	"time"
)

func (s *Scraper) fixFileReferences(url *url.URL, buf io.Reader) (string, error) {
	g, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return "", err
	}

	relativeToRoot := s.urlRelativeToRoot(url)

	g.Find("a").Each(func(_ int, selection *goquery.Selection) {
		s.fixQuerySelectionForPattern(url, "href", selection, true, relativeToRoot)
	})

	g.Find("link").Each(func(_ int, selection *goquery.Selection) {
		s.fixQuerySelection(url, "href", selection, false, relativeToRoot)
	})

	g.Find("img").Each(func(_ int, selection *goquery.Selection) {
		s.fixQuerySelection(url, "src", selection, false, relativeToRoot)
	})

	g.Find("script").Each(func(_ int, selection *goquery.Selection) {
		s.fixQuerySelection(url, "src", selection, false, relativeToRoot)
	})

	return g.Html()
}

func (s *Scraper) fixQuerySelection(url *url.URL, attribute string, selection *goquery.Selection,
	linkIsAPage bool, relativeToRoot string) {
	src, ok := selection.Attr(attribute)
	if !ok {
		return
	}

	if strings.HasPrefix(src, "data:") {
		return
	}
	if strings.HasPrefix(src, "mailto:") {
		return
	}
	resolved := s.resolveURL(url, src, linkIsAPage, relativeToRoot)
	if src == resolved { // nothing changed
		return
	}

	s.log.Debug("HTML Element relinked", zap.String("URL", src), zap.String("Fixed", resolved))
	selection.SetAttr(attribute, resolved)
}

func (s *Scraper) fixQuerySelectionForPattern(url *url.URL, attribute string, selection *goquery.Selection,
	linkIsAPage bool,relativeToRoot string){


	src, ok := selection.Attr(attribute)
	if !ok {
		return
	}


	tmpSrc:=src
	finalPath:=""

	println("src is: ",src)
	var hostPresent bool
	hostPresent=true

	//if strings.Contains(tmpSrc,"http://") || strings.Contains(tmpSrc,"https://") {
	if tmpSrc!="#" && tmpSrc!="/"{
		if !strings.Contains(tmpSrc,url.Host){
			hostPresent=false
			tmpSrc=url.Host+"/"+tmpSrc
		}

		protocol:=""
		if strings.Contains(tmpSrc,"https://"){
			protocol="https://"
			tmpSrc=strings.TrimPrefix(tmpSrc,"https://")
		}else if strings.Contains(tmpSrc,"http://"){
			tmpSrc=strings.TrimPrefix(tmpSrc,"http://")
			protocol="http://"
		}

		//if strings.Contains(tmpSrc,"pricing-options"){
		//	println("pricing option lai milaudai chu hai guys")
		//	time.Sleep(5*time.Second)
		//}

		if strings.Contains(tmpSrc,".html"){
			if !strings.Contains(tmpSrc,"index.html"){
				println("doesn't contain index.html but has other .html page hai")
				time.Sleep(2*time.Second)
			}
		}

		pathArr:=strings.Split(tmpSrc,"/")
		newPathArr:=[]string{}
		for _,val:=range pathArr{
			if val!=""{
				newPathArr=append(newPathArr,val)
			}
		}

		length:=len(newPathArr)
		for i:=0;i<length-1;i++{
				finalPath+= newPathArr[i] + "/"
		}
		println("new path arr is : ",newPathArr)

		if val,ok:=s.Config.FolderCount[finalPath];ok && length!=0{
			println("match coming inside for: ",finalPath)
			println("newpath len-1 is: ",newPathArr[length-1])
			if _,exists:=val[newPathArr[length-1]];!exists{
				println("doesn't exist and needs to resolve: ",newPathArr[length-1])

				//provide an already existing one as url
				path:=strings.TrimPrefix(finalPath,url.Host)

				for k,_:=range val{
					if pathArr[len(pathArr)-1]==""{
						path+=k +"/"
					}else{
						path+=k
					}
					break
				}
				if hostPresent{
					tmpSrc=protocol+url.Host+path
				}else{
					tmpSrc=path
				}
				src=tmpSrc
			}
		}
	}

	fmt.Println("src changed into : ",src)

	resolved := s.resolveURL(url, src, linkIsAPage, relativeToRoot)
	if src == resolved { // nothing changed
		return
	}

	fmt.Println("resolved into: ",resolved)
	s.log.Debug("HTML Element relinked", zap.String("URL", src), zap.String("Fixed", resolved))
	selection.SetAttr(attribute, resolved)
}