package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Scrape(c *gin.Context) {
	log.Println("[[handler/scrape]] New HTTP Request")

	/*
	   TODO:
	   1. Take in request body url params
	   2. Start scrapper
	   3. Respond with scrape_id and status 200
	   Sample POST request body
	   {
	       "url":"www.airbnb.com/hosting",
	       "screen_width":1920,
	       "screen_height": 1080,
	       "username": "abc@def.gh",
	       "password": "L36gh!h'",
	        "project_id": "abc", //optional
	       "folder_threshold": 20,
	        "folder_examples_count":3,
	       "patterns": ["www.airbnb.com/s/asterisk(*)/experiences]"

	*/
	response:=make(map[string]interface{})
	response["job_id"]="12345678"
	response["msg"]="Scrapping Started"

	c.JSON(
		http.StatusOK,
		response,
	)

}
