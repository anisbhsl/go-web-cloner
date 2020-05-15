package scraper

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
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

	/*
		URL : url => columbus-internet.com/de/technologien/

	*/
	src, ok := selection.Attr(attribute)
	if !ok {
		return
	}


	tmpSrc:=src
	finalPath:=""

	println("url is: ",url.Host+url.Path)
	println("src is: ",src)

	if strings.Contains(tmpSrc,"http://") || strings.Contains(tmpSrc,"https://") {
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

		if val,ok:=s.Config.FolderCount[finalPath];ok && length!=0{
			println("match coming inside for: ",finalPath)
			println("newpath len-1 is: ",newPathArr[length-1])
			if _,exists:=val[newPathArr[length-1]];!exists{
				println("doesn't exist and needs to resolve: ",newPathArr[length-1])

				//provide an already existing one as url
				path:=strings.TrimPrefix(finalPath,url.Host)

				for k,_:=range val{
					//if strings.Contains(k,".html"){
					//	path+=k
					//}else{
						path+=k+"/"
					//}
					break
				}
				tmpSrc=protocol+url.Host+"/"+path
				src=tmpSrc
			}
		}
	}

	resolved := s.resolveURL(url, src, linkIsAPage, relativeToRoot)
	if src == resolved { // nothing changed
		return
	}

	s.log.Debug("HTML Element relinked", zap.String("URL", src), zap.String("Fixed", resolved))
	selection.SetAttr(attribute, resolved)



}