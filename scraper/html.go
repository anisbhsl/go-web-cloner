package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"io"
	"log"
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

	//should fix issue for playheartstone
	//if strings.HasPrefix(src,"//"){
	//	src=s.URL.Scheme+src
	//}

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

	doesNotExist:=false

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
		println("finalpath is: ",finalPath)
		time.Sleep(2*time.Second)

		if val,ok:=s.Config.FolderCount[finalPath];ok && length!=0{
			println("match coming inside for: ",finalPath)
			println("newpath len-1 is: ",newPathArr[length-1])
			if _,exists:=val[newPathArr[length-1]];!exists{
				println("doesn't exist and needs to resolve: ",newPathArr[length-1])
				println("val is : %v",val)

				//provide an already existing one as url
				//path:=strings.TrimPrefix(finalPath,url.Host)
				//
				//for k,_:=range val{
				//	if pathArr[len(pathArr)-1]==""{
				//		path+=k +"/"
				//	}else{
				//		path+=k
				//	}
				//	break
				//}
				//if hostPresent{
				//	tmpSrc=protocol+url.Host+path
				//}else{
				//	tmpSrc=path
				//}
				//src=tmpSrc
				//
				//var srcExists bool
				//for _,val:=range s.report.DetailedReport{
				//	if src==val.OriginURL==src{
				//		srcExists=true
				//	}
				//}

				loopCount:=0
				srcExists:=false

				path:=strings.TrimPrefix(finalPath,url.Host)
				for k,_:=range val{
					tmpPath:=path
					//if pathArr[len(pathArr)-1]==""{  //todo: needs a good workaround here hai
					//	path+=k +"/"
					//}else{
					//	path+=k
					//}
					path+=k

					if hostPresent{
						tmpSrc=protocol+url.Host+path
					}else{
						tmpSrc=path
					}
					src=tmpSrc

					for _,v:=range s.report.DetailedReport{
						if src==v.OriginURL && v.StatusCode<400{
							srcExists=true
							break
						}
					}

					if srcExists==true || loopCount>=5{
						break
					}
					path=tmpPath
					loopCount++
					println("updated path: ",path)
				}
				doesNotExist=true
			}
		}else{
			finalPath:=url.Host+"/"
			//for host part
			if _,ok:=s.Config.FolderCount[finalPath];!ok{
				folder:=make(map[string]bool)
				folder[newPathArr[0]]=true

				s.Config.FolderCount[finalPath]=folder
			}

			length:=len(newPathArr)
			for i:=0;i<length-1;i++ {

				finalPath += newPathArr[i]
				finalPath += "/"

				if _, ok := s.Config.FolderCount[finalPath]; !ok {
					folder := make(map[string]bool)

					folder[newPathArr[i+1]] = true
					s.Config.FolderCount[finalPath] = folder

				}
			}


		}
	}


	fmt.Println("src changed into : ",src)

	resolved := s.resolveURL(url, src, linkIsAPage, relativeToRoot)
	if src == resolved { // nothing changed
		return
	}

	if doesNotExist{
		log.Println("doesn't exist and added to mandatory download")
		s.Config.MandatoryDownload[src]=true
	}

	fmt.Println("resolved into: ",resolved)
	s.log.Debug("HTML Element relinked", zap.String("URL", src), zap.String("Fixed", resolved))
	selection.SetAttr(attribute, resolved)
}