package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

//Redirect handler redirects the HTTP request to destination URL
//This is used in reports
//Requests come at endpoint: /api/redirect?url=<url>
func Redirect(c *gin.Context){
	//url:=c.Param("url")
	url:=c.Query("url")
	log.Println("[[INFO]] Redirect request to: ",url)
	c.Redirect(307,"http://"+url)
	return
}