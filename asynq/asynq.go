package asyncq

import (
	"go-web-cloner/scraper"
)

type Dispatcher struct{
	Queue []string
	Scraper *scraper.Scraper
}

func NewDispatcher()*Dispatcher{
	return &Dispatcher{
			Queue:[]string{},
		}
}

func (d *Dispatcher) IsWorkerAvailable() bool{
	if len(d.Queue)>=1{
		return false
	}
	//available
	return true
}

func (d *Dispatcher) StopScrapper(){
	d.Scraper.Config.Stop=true
	d.Queue=[]string{}
}