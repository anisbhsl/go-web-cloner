package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"io"
	"net/url"
	"strings"
)

func (s *Scraper) fixFileReferences(url *url.URL, buf io.Reader) (string, error) {
	g, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return "", err
	}

	relativeToRoot := s.urlRelativeToRoot(url)
	//fmt.Println("url: ",url.Path)
	//fmt.Println("relativeToRoot URL: ",relativeToRoot)

	//get list of browser links to fix
	//var linksToFix []*browser.Link

	//holds path and folders inside that path
	//e.g.: folderCount["team"]=["john","doe","david"]
//	folderCount:=make(map[string]map[string]bool)
/*
	for _,link:=range s.browser.Links(){
		linksToFix=append(linksToFix,link)
		println("link: ",link.URL.Host+link.URL.Path)
		//check which links exceed folder count
		path:=link.URL.Path
		if path=="" || path=="/"{
			continue
		}
		pathArr:=strings.Split(path,"/")

		newPathArr:=[]string{} //will hold path array only
		for _,val:=range pathArr{
			if val!=""{
				newPathArr=append(newPathArr,val)
			}
		}

		finalPath:=link.URL.Host+"/"
		length:=len(newPathArr)
		for i:=0;i<length-1;i++{
			finalPath+=newPathArr[i]+"/"
		}

		if length==0{
			continue
		}

		if val,ok:=s.Config.FolderCount[finalPath];!ok{
			folder:=make(map[string]bool)
			folder[newPathArr[length-1]]=true
			folderCount[finalPath]=folder
		}else{
			if len(val)>=s.Config.FolderThreshold{
				continue
			}else{
				if newPathArr[length-1]=="pricing-options"{
					time.Sleep(5*time.Second)
					println("newpatharr ma pricing option aayo")
				}
				val[newPathArr[length-1]]=true
				folderCount[finalPath]=val
			}
		}
	}
*/
	/*
		k,v====> columbus-internet.com/ 0xc00054efc0
		key:  de
		k,v====> columbus-internet.com/de/ 0xc00054ef30
		key:  technologien
		key:  preisegestaltung-optionen
		key:  perfekte-software
		key:  cookie-richtlinie
		k,v====> columbus-internet.com/en/ 0xc00054ef60
		key:  technologies
		key:  legal-notice
		key:  cookie-policy
		k,v====> columbus-internet.com/de/unternehmen/ 0xc00054ef90
		key:  referenzprojekte
		key:  team
	*/

	//println("foldercount: ........")
	//for k,v:=range folderCount{
	//	println("k,v====>",k,v)
	//	for key,_:=range v{
	//		println("key: ",key)
	//	}
	//}

	g.Find("a").Each(func(_ int, selection *goquery.Selection) {
		s.fixQuerySelectionWithPattern(url, "href", selection, true, relativeToRoot)
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

//fixQuerySelectionWithPattern fixes links which point to page not scraped
//due to the folder_threshold and makes them point to one of
//the matching pages which was scraped
func (s *Scraper) fixQuerySelectionWithPattern(url *url.URL, attribute string, selection *goquery.Selection,
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

	/*
	URL : url => columbus-internet.com/de/technologien/

	*/

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
			finalPath+=newPathArr[i]+"/"
		}

		if val,ok:=s.Config.FolderCount[finalPath];ok && length!=0{
			println("match coming inside for: ",finalPath)
			println("newpath len-1 is: ",newPathArr[length-1])
			if _,exists:=val[newPathArr[length-1]];!exists{
				println("doesn't exist and needs to resolve: ",newPathArr[length-1])

				//provide an already existing one as url
				path:=strings.TrimPrefix(finalPath,url.Host)

				for k,_:=range val{
					path+=k+"/"
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
	println("resolved url is : ",resolved)

	/*
		resolved url:  ../../team/evgeniy-b/index.html
		resolved url:  ../../team/anita-s/index.html
		resolved url:  ../../team/anita-s/index.html
		resolved url:  ../../team/andrei-a/index.html
		resolved url:  ../../team/andrei-a/index.html

		.html cha ki nai hernu paryo... ani tyo relativeToRoot kasto audo raicha hernu paryo

		relativeToRoot:  ../../../
		resolved url:  ../../team/andrew-sittermann/index.html


	*/
	//println("-------------START-----------------")
	//println("relativeToRoot: ",relativeToRoot)
	//println("resolved url: ",resolved)
	//println("--------------END------------------")
	s.log.Debug("HTML Element relinked", zap.String("URL", src), zap.String("Fixed", resolved))
	selection.SetAttr(attribute, resolved)
}



