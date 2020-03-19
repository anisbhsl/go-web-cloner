package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func Scrape() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
}
